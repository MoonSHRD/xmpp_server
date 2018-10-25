package xep0045

import (
    "github.com/ethereum/go-ethereum/accounts/keystore"
    "github.com/ortuman/jackal/helpers"
    "github.com/ortuman/jackal/model"
    "github.com/ortuman/jackal/module/xep0030"
    "github.com/ortuman/jackal/router"
    "github.com/ortuman/jackal/storage"
    "github.com/ortuman/jackal/stream"
    "github.com/ortuman/jackal/xml"
    "github.com/ortuman/jackal/xml/jid"
	"log"
	"sort"
    "strconv"
    "strings"
)

const chatNamespace = "http://jabber.org/protocol/muc"
const chatEventNamespace = "http://jabber.org/protocol/muc#event"
const discoNamespace = "http://jabber.org/protocol/disco#items"

type RegisterChat struct {
    stm        stream.C2S
}

type user_role struct {
    affiliation string
    role string
}

var roles = struct {
    owner user_role
    admin user_role
    paticipant user_role
}{
    owner:user_role{"owner","moderator"},
    admin:user_role{"admin","moderator"},
    paticipant:user_role{"member","paticipant"},
}

type chat_type struct {
    kind string
}

var types = struct {
    user      chat_type
    group      chat_type
    channel chat_type
}{
    user:      chat_type{"user_chat"},
    group:      chat_type{"group"},
    channel: chat_type{"channel"},
}

func (x *RegisterChat) RegisterDisco(discoInfo *xep0030.DiscoInfo) {
    // register disco feature
    discoInfo.Entity(x.stm.Domain(), "").AddFeature(chatNamespace)
}

func generateRoleItem(role user_role) xml.XElement {
    item_elem:=xml.NewElementName("item")
    item_elem.SetAttribute("affiliation",role.affiliation)
    item_elem.SetAttribute("role",role.role)
    return item_elem
}


func New(stm stream.C2S) *RegisterChat {
    return &RegisterChat{
        stm: stm,
    }
}

func sendChatEvent(user model.User,)  {

}


func (x *RegisterChat) CreateChat(presence *xml.Presence) {
    var err error
    to:=presence.ToJID()
    from:=presence.FromJID()
    kind := types.channel.kind
    if presence.Attributes().Get("channel")=="2" {
        kind = types.group.kind
    }
    chat := model.Chat{Chatname:to.Node(),Type:kind,Creator:from.Node()}
    //todo: deal with double chat insert
    newAcc, err := keystore.NewKeyStore("", keystore.StandardScryptN, keystore.StandardScryptP).NewAccount(from.Node())
    chat.Id = strings.ToLower(newAcc.Address.Hex())
    seed, _ := strconv.ParseInt(chat.Id, 10, 64)
    chat.Avatar = helpers.GenerateThumb(seed)
    chat.Id, err = storage.Instance().InsertOrUpdateChat(&chat)
    storage.Instance().InsertChatUser(chat.Id,from.Node(), roles.owner.affiliation)
    if err != nil {
        x.stm.SendElement(presence.NotAllowedError())
    } else {
        x.sendJoinAcceptance(from, &chat, roles.owner)
    }
}


func (x *RegisterChat) sendJoinEvent(chat_id string,user *jid.JID, date string) {
    chat,_:=storage.Instance().FetchChat(chat_id)
    x.sendJoinAcceptance(user,chat,roles.paticipant)
//todo: Добавить отправление времени сообщений
    x_elem:=xml.NewElementName("x")
    x_elem.SetNamespace(chatNamespace+"#user")

    elem:=xml.NewElementName("presence")
    elem.SetFrom(chat_id + "@localhost")
    elem.SetAttribute("user_joined",user.NDString())
    elem.SetAttribute("date", date)
    elem.AppendElement(x_elem)

    chat_u,_ := storage.Instance().FetchChatUsers(chat_id)

    x.sendToUsers(elem,chat_u)
}

func (x *RegisterChat) sendJoinAcceptance(user *jid.JID,chat *model.Chat,role user_role) {
	elem:=xml.NewElementName("presence")
	elem.SetAttribute("channel",string(chat.Type))
	elem.SetAttribute("avatar",chat.Avatar)
	elem.SetFrom(chat.Id + "@localhost/"+chat.Chatname)
	elem.SetTo(user.NDString())
    s_elem := xml.NewElementName("set")
    messages, _ := storage.Instance().GetMsgsFromDB(chat.Id)
    for _, message := range(messages) {
        item := xml.NewElementName("item")
        item.SetAttribute("sender", message.Sender)
        item.SetAttribute("message", message.Message)
		item.SetAttribute("time", message.Time.String())

		if (message.File != 0) {
			files, err := storage.Instance().GetFilesFromDB(message.Id)

			if (err != nil) {
				return
			}

			file_elem := xml.NewElementName("files")
			for _, file:=range(files) {
				item := xml.NewElementName("item")
				item.SetAttribute("name", file.Name)
				item.SetAttribute("type", file.Type)
				item.SetText(file.Hash)
				file_elem.AppendElement(item)
			}
			item.AppendElement(file_elem)
		}

        s_elem.AppendElement(item)
    }
    x_elem:=xml.NewElementName("x")
    x_elem.SetNamespace(chatNamespace+"#user")
    x_elem.AppendElement(generateRoleItem(role))


    elem.AppendElement(x_elem)
    elem.AppendElement(s_elem)

    x.stm.SendElement(elem)
}


func (x *RegisterChat) ProcessPresence(presence *xml.Presence) {
    to:=presence.ToJID()
    from:=presence.FromJID()
    groupName := to.Node()

    exist, err:=storage.Instance().ChatExists(groupName)
    if !exist || err!=nil {
        x.CreateChat(presence)
        return
    }
    //todo Защита от перезаписи админа
    date, _ := storage.Instance().InsertChatUser(groupName, from.Node(),roles.paticipant.affiliation)
    x.sendJoinEvent(groupName,from, date)
}

func (x *RegisterChat) sendToUsers(elem *xml.Element, users model.ChatUsers) {
    for username,_ := range users {
        elem.SetTo(username)
        for _,u_stream := range router.UserStreams(username) {
            u_stream.SendElement(elem)
        }
    }
}

func (x *RegisterChat) sendToAnotherUser(elem *xml.Element, username string) {
    for _,u_stream := range router.UserStreams(username) {
        u_stream.SendElement(elem)
    }
}

func (x *RegisterChat) ProcessMessage(msg *xml.Message) {

    node := msg.ToJID().Node()
    if node == "" {
        x.stm.SendElement(msg.BadRequestError())
        return
    }

    id:=node

    exist,err := storage.Instance().ChatExists(id)

    if err != nil {
        x.stm.SendElement(msg.BadRequestError())
        return
    }

    if !exist {
        x.stm.SendElement(msg.BadRequestError())
        return
    }

    chat,_ := storage.Instance().FetchChat(id)
    chat_u,_ := storage.Instance().FetchChatUsers(id)
    if chat.Type == "channel" && chat_u[msg.FromJID().Node()].Role!="owner" {
        x.stm.SendElement(msg.BadRequestError())
        return
    }

    if chat.Type == "group" {
        msg.SetAttribute("sender",msg.From())
    }
    body := (xml.NewElementName("body"))
    body.SetText(msg.Elements().Child("body").Text())

    elem := xml.NewElementName("message")
    elem.SetAttribute("sender", msg.To())
    elem.SetAttribute("from", msg.To())
    elem.SetAttribute("type", "groupchat")
    elem.AppendElement(body)

    message := msg.Elements().Child("body")
    id_user := msg.Elements().Child("id")
    list_files := x.GetFilesFromStanza(msg)
	with_file := 0
	if len(list_files) > 0 {
		with_file = 1

		file_elem := xml.NewElementName("files")
		for _, file:=range(list_files) {
			item := xml.NewElementName("item")
			item.SetAttribute("name", file.Name)
			item.SetAttribute("type", file.Type)
			item.SetText(file.Hash)
			file_elem.AppendElement(item)
		}
		elem.AppendElement(file_elem)

	}
    id_db, date, err := storage.Instance().WriteMsgToDB(id, id, message.Text(), 1, with_file)
    if err != nil {
    	log.Println(err)
		x.stm.SendElement(msg.BadRequestError())
    	return
	}

	// Записываем все файлы из станзы в базу
	for _, file:=range(list_files) {
		err := storage.Instance().WriteFileToDB(file, id_db)
		if err != nil {
			x.stm.SendElement(msg.BadRequestError())
			return
		}
	}

    x.SendConfirmation(id_user, int(id_db), msg.FromJID().Node(), date)

    x_elem := xml.NewElementName("x")
    x_elem.SetAttribute("stamp", date)
    elem.AppendElement(x_elem)
    delete(chat_u, msg.FromJID().Node())

    x.sendToUsers(elem,chat_u)
}

func (x *RegisterChat) ProcessElem(stanza xml.Stanza) (string, bool) {

    from:=stanza.FromJID()
    to:=stanza.ToJID()
    ok:=to!=nil && from != nil

    if !ok {
        return "", false
    }

    switch stanza:=stanza.(type) {
    case *xml.Presence:

        el:=stanza.Elements().ChildNamespace("x", chatNamespace)
        if el == nil{
            return "", false
        }
        x.ProcessPresence(stanza)

    case *xml.Message:

        if stanza.Type() == "chat" {
            is_user, err := storage.Instance().UserExists(to.Node())
            if is_user == false || err != nil {
                return "", true
            }
            stanz_elems := stanza.Element.Elements()
            msg := stanz_elems.Child("body")
            id_user := stanz_elems.Child("id")
            users_id := []string{from.Node(), to.Node()}
            sort.Strings(users_id)
            chat_id := strings.Join(users_id, "_")
            exist, _ := storage.Instance().ChatExists(chat_id)
            if !exist {
                chat := model.Chat{Id:chat_id, Chatname:"", Type:"user_chat", Creator:"server"}
                storage.Instance().InsertOrUpdateChat(&chat)
                storage.Instance().InsertChatUser(chat_id, users_id[0], "")
                storage.Instance().InsertChatUser(chat_id, users_id[1], "")
            }
            list_files := x.GetFilesFromStanza(stanza)
			with_file := 0
			if len(list_files) > 0 {
				with_file = 1
			}
            id_db, date, err := storage.Instance().WriteMsgToDB(chat_id, from.Node(), msg.Text(), 1, with_file)
            if err != nil {
                return "", true
            }
            // Записываем все файлы из станзы в базу
			for _, file:=range(list_files) {
				err := storage.Instance().WriteFileToDB(file, id_db)
				if err != nil {
					return "", true
				}
			}

			x.SendConfirmation(id_user, int(id_db), stanza.FromJID().NDString(), date)
            return date, false
        }
        x.ProcessMessage(stanza)

    case *xml.IQ:
        if stanza.Elements().ChildNamespace("query", discoNamespace)!= nil{
            x.FindGroup(stanza)
            return "", true
        }
        if stanza.Elements().ChildNamespace("x", chatEventNamespace)!= nil{
            x.ProcessChatEvent(stanza)
            return "", true
        }
        return "", false
    }

    return "", true
}

func (x *RegisterChat) FindGroup(presence *xml.IQ){
    q_elem := xml.NewElementName("query")
    q_elem.SetNamespace(discoNamespace)
    a := presence.Attributes().Get("name")
    res := storage.Instance().FindGroups(a)
    for _, group := range(res){
       item := xml.NewElementName("item")
       item.SetAttribute("jid", group.Id + "@localhost")
       item.SetAttribute("name", group.Chatname)
       item.SetAttribute("avatar", group.Avatar)
       item.SetAttribute("channel", string(group.Type))
       q_elem.AppendElement(item)
    }
    elem := xml.NewElementName("iq")
    elem.SetFrom("localhost")
    elem.SetTo(presence.FromJID().NDString())
    elem.SetType("result")
    elem.AppendElement(q_elem)
    x.stm.SendElement(elem)
}

func (x *RegisterChat) ProcessChatEvent(iq *xml.IQ){

    chat_id := iq.ToJID().Node()
    chat,err:=storage.Instance().FetchChat(chat_id)
    if err !=nil || chat==nil {
        x.stm.SendElement(iq.BadRequestError())
        return
    }

    eventItem:=iq.Elements().Child("x").Elements().Child("item")
    switch eventItem.Attributes().Get("type") {
    case "suggestion":
        elem:=xml.NewElementFromElement(iq)
        elem.SetTo(chat.Creator+"@localhost")
        x.sendToAnotherUser(elem,chat.Creator)
    }
}
func (x *RegisterChat) SendConfirmation(idUser xml.XElement, id_db int, to string, date string) {
    q_elem := xml.NewElementName("confirmation")
    q_elem.SetNamespace(discoNamespace)
    item := xml.NewElementName("item")
    item.SetAttribute("userid", "1")
    item.SetAttribute("DBid", strconv.Itoa(int(id_db)))
    item.SetAttribute("date", date)
    q_elem.AppendElement(item)
    elem := xml.NewElementName("iq")
    elem.SetFrom("localhost")
    elem.SetTo(to)
    elem.SetType("result")
    elem.AppendElement(q_elem)
    x.stm.SendElement(elem)
}

func (x *RegisterChat) GetFilesFromStanza(msg *xml.Message) []model.File{
	attrs := msg.Elements()
	var list_files []model.File
	is_files := attrs.Child("files")
	if is_files != nil {
		for _, file := range (is_files.Elements().All()) {
			if file != nil {
				hash := file.Text()
				file_type := file.Attributes().Get("type")
				name := file.Attributes().Get("name")
				_file := model.File{Hash: hash, Type: file_type, Name: name}
				list_files = append(list_files, _file)
			}
		}
	}
	return list_files
}
