package xep0045

import (
    "github.com/ortuman/jackal/xml"
    "github.com/ortuman/jackal/stream"
    "github.com/ortuman/jackal/module/xep0030"
    "github.com/ortuman/jackal/storage"
    "github.com/ortuman/jackal/model"
    "strconv"
    "github.com/ortuman/jackal/xml/jid"
    "github.com/ortuman/jackal/router"
    "time"
    "github.com/ortuman/jackal/helpers"
    "github.com/ethereum/go-ethereum/accounts/keystore"
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
    //contractaddress := presence.Attributes().Get("contractaddress")
    channel := false
    if presence.Attributes().Get("channel")=="1" {
        channel = true
    }
    chat:=model.Chat{Chatname:to.Node(),Channel:channel,Creator:from.Node()}
    //todo: deal with double chat insert
    newAcc, err := keystore.NewKeyStore("", keystore.StandardScryptN, keystore.StandardScryptP).NewAccount(from.Node())
    chat.Id = newAcc.Address.Hex()
    seed, _ := strconv.ParseInt(chat.Id, 10, 64)
    chat.Avatar = helpers.GenerateThumb(seed)
    chat.Id, err = storage.Instance().InsertOrUpdateChat(&chat)
    storage.Instance().InsertChatUser(chat.Id,from.Node(),true)
    if err != nil {
        x.stm.SendElement(presence.NotAllowedError())
    } else {
        x.sendJoinAcceptance(from, &chat, roles.owner)
        //x_elem:=xml.NewElementName("x")
        //x_elem.SetNamespace(chatNamespace+"#user")
        //x_elem.AppendElement(generateRoleItem(roles.owner))
        //
        //elem:=xml.NewElementName("presence")
        //elem.SetFrom(strconv.Itoa(int(chat.Id))+"@localhost/"+to.Node())
        //elem.SetTo(from.NDString())
        //elem.SetAttribute("channel",chat.IsChannel())
        //elem.AppendElement(x_elem)
        //
        //x.stm.SendElement(elem)
    }
}


func (x *RegisterChat) sendJoinEvent(chat_id string,user *jid.JID) {
    chat,_:=storage.Instance().FetchChat(chat_id)
    x.sendJoinAcceptance(user,chat,roles.paticipant)
    
    x_elem:=xml.NewElementName("x")
    x_elem.SetNamespace(chatNamespace+"#user")
    
    elem:=xml.NewElementName("presence")
    elem.SetFrom(chat_id + "@localhost")
    elem.SetAttribute("user_joined",user.NDString())
    //elem.SetAttribute("contractaddress", contractaddress)
    elem.AppendElement(x_elem)
    
    chat_u,_ := storage.Instance().FetchChatUsers(chat_id)
    
    x.sendToUsers(elem,chat_u)
    //for username,_ := range chat_u {
    //    elem.SetTo(username)
    //    for _,u_stream := range router.UserStreams(username) {
    //        u_stream.SendElement(elem)
    //    }
    //}
}

func (x *RegisterChat) sendJoinAcceptance(user *jid.JID,chat *model.Chat,role user_role) {
    
    x_elem:=xml.NewElementName("x")
    x_elem.SetNamespace(chatNamespace+"#user")
    x_elem.AppendElement(generateRoleItem(role))
    
    elem:=xml.NewElementName("presence")
    elem.SetAttribute("channel",chat.IsChannel())
    //elem.SetAttribute("contractaddress",contractaddress)
    elem.SetAttribute("avatar",chat.Avatar)
    elem.SetFrom(chat.Id + "@localhost/"+chat.Chatname)
    elem.SetTo(user.NDString())
    elem.AppendElement(x_elem)
    
    x.stm.SendElement(elem)
}


func (x *RegisterChat) ProcessPresence(presence *xml.Presence) {
    //var err error
    to:=presence.ToJID()
    from:=presence.FromJID()
    //contractaddress := presence.Attributes().Get("contractaddress")
    groupName := to.Node()
    //x.CreateChat(presence)
    //if err!=nil{
    //    x.CreateChat(presence)
    //    return
    //}
    exist,err:=storage.Instance().ChatExists(groupName)
    if !exist || err!=nil {
        x.CreateChat(presence)
        return
    }
    //todo Защита от перезаписи админа
    storage.Instance().InsertChatUser(groupName, from.Node(),false)
    x.sendJoinEvent(groupName,from)
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
    //if err != nil {
    //    x.stm.SendElement(msg.BadRequestError())
    //    return
    //}
    
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
    
    if chat.Channel && chat_u[msg.FromJID().Node()].Admin!=1 {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    if chat.Channel {
        msg.SetAttribute("sender",msg.From())
    }
    
    elem:=xml.NewElementFromElement(msg)
    elem.SetFrom(msg.To())
    x_elem:=xml.NewElementName("x")
    x_elem.SetAttribute("stamp",time.Now().String())
    msg.AppendElement(x_elem)
    
    x.sendToUsers(elem,chat_u)
    
    //for username,_ := range chat_u {
    //    msg.SetTo(username)
    //    x.stm.SendElement(msg)
    //}
}

func (x *RegisterChat) ProcessElem(stanza xml.Stanza) bool {
    
    from:=stanza.FromJID()
    to:=stanza.ToJID()
    ok:=to!=nil && from != nil
    
    if !ok {
        return false
    }
    
    switch stanza:=stanza.(type) {
    case *xml.Presence:
        
        el:=stanza.Elements().ChildNamespace("x", chatNamespace)
        if el == nil{
            return false
        }
        x.ProcessPresence(stanza)

    case *xml.Message:
    
        if !stanza.IsGroupChat() {
            stanz_elems := stanza.Element.Elements()
            msg := stanz_elems.Child("body")
            id_user := stanz_elems.Child("id")
            for _, user := range router.UserStreams(to.Node()) {
                if user.Username() == to.Node() {
                    id_db, err := storage.Instance().WriteMsgToDB(to.Node(), from.Node(), msg.Text(), 1)
                    if err != nil {
                        return false
                    }
                    q_elem := xml.NewElementName("confirmation")
                    q_elem.SetNamespace(discoNamespace)
                    item := xml.NewElementName("item")
                    item.SetAttribute("userid", id_user.Text())
                    item.SetAttribute("DBid", strconv.Itoa(int(id_db)))
                    q_elem.AppendElement(item)
                    elem := xml.NewElementName("iq")
                    elem.SetFrom("localhost")
                    elem.SetTo(stanza.FromJID().NDString())
                    elem.SetType("result")
                    elem.AppendElement(q_elem)
                    x.stm.SendElement(elem)
                    return false
                }
            }
            return false
        }
        x.ProcessMessage(stanza)

    case *xml.IQ:
        if stanza.Elements().ChildNamespace("query", discoNamespace)!= nil{
            x.FindGroup(stanza)
            return true
        }
        if stanza.Elements().ChildNamespace("x", chatEventNamespace)!= nil{
            x.ProcessChatEvent(stanza)
            return true
        }
        return false
    }

    return true
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
       //item.SetAttribute("contractaddress", group.Contractaddress)
       item.SetAttribute("avatar", group.Avatar)
       item.SetAttribute("channel", group.IsChannel())
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
    //if err !=nil {
    //    x.stm.SendElement(iq.BadRequestError())
    //    return
    //}
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
        //x.stm.SendElement(elem)
    }
}
