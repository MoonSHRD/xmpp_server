package xep0045

import (
    "github.com/ortuman/jackal/xml"
    "github.com/ortuman/jackal/stream"
    "github.com/ortuman/jackal/module/xep0030"
    "github.com/ortuman/jackal/storage"
    "github.com/ortuman/jackal/model"
    "strconv"
    "time"
    "github.com/ortuman/jackal/xml/jid"
)

const chatNamespace = "http://jabber.org/protocol/muc"

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


func (x *RegisterChat) CreateChat(presence *xml.Presence) {
    var err error
    to:=presence.ToJID()
    from:=presence.FromJID().NDString()
    chat:=model.Chat{Chatname:to.Node(),Channel:false,Creator:from}
    chat.Id, err = storage.Instance().InsertOrUpdateChat(&chat)
    storage.Instance().InsertChatUser(chat.Id,from,true)
    if err != nil {
        x.stm.SendElement(presence.NotAllowedError())
    } else {
        x_elem:=xml.NewElementName("x")
        x_elem.SetNamespace(chatNamespace+"#user")
        x_elem.AppendElement(generateRoleItem(roles.owner))
        
        elem:=xml.NewElementName("presence")
        elem.SetFrom(strconv.Itoa(int(chat.Id))+"@localhost/"+to.Node())
        elem.SetTo(from)
        elem.AppendElement(x_elem)
        
        x.stm.SendElement(elem)
    }
}


func (x *RegisterChat) sendJoinEvent(chat_id int64,user *jid.JID) {
    chat,_:=storage.Instance().FetchChat(chat_id)
    x.sendJoinAcceptance(user,chat,roles.paticipant)
    
    x_elem:=xml.NewElementName("x")
    x_elem.SetNamespace(chatNamespace+"#user")
    
    elem:=xml.NewElementName("presence")
    elem.SetFrom(strconv.FormatInt(chat_id,10)+"@localhost")
    elem.SetAttribute("user_joined",user.NDString())
    elem.AppendElement(x_elem)
    
    
    chat_u,_ := storage.Instance().FetchChatUsers(chat_id)
    for _,username := range chat_u {
        elem.SetTo(username)
        x.stm.SendElement(elem)
    }
}

func (x *RegisterChat) sendJoinAcceptance(user *jid.JID,chat *model.Chat,role user_role) {
    
    x_elem:=xml.NewElementName("x")
    x_elem.SetNamespace(chatNamespace+"#user")
    x_elem.AppendElement(generateRoleItem(role))
    
    elem:=xml.NewElementName("presence")
    elem.SetFrom(strconv.Itoa(int(chat.Id))+"@localhost/"+chat.Chatname)
    elem.SetTo(user.NDString())
    elem.AppendElement(x_elem)
    
    x.stm.SendElement(elem)
}


func (x *RegisterChat) ProcessPresence(presence *xml.Presence) {
    //var err error
    to:=presence.ToJID()
    from:=presence.FromJID()
    id,err:=strconv.ParseInt(to.Node(),10,64)
    if err!=nil{
        x.CreateChat(presence)
        return
    }
    exist,err:=storage.Instance().ChatExists(id)
    if exist && err!=nil {
        x.CreateChat(presence)
        return
    }
    
    storage.Instance().InsertChatUser(id,from.NDString(),false)
    x.sendJoinEvent(id,from)
}

func (x *RegisterChat) ProcessMessage(msg *xml.Message) {
    
    node := msg.ToJID().Node()
    if node == "" {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    id,err:=strconv.ParseInt(node,10,64)
    if err != nil {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    exist,err := storage.Instance().ChatExists(id)
    
    if err != nil {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    if !exist {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    chat_u,_ := storage.Instance().FetchChatUsers(id)
    //chat,_ := storage.Instance().FetchChat(id)
    
    msg.SetAttribute("sender",msg.From())
    msg.SetFrom(msg.To())
    for _,username := range chat_u {
        msg.SetTo(username)
        x_elem:=xml.NewElementName("x")
        x_elem.SetAttribute("stamp",time.Now().String())
        msg.AppendElement(x_elem)
        x.stm.SendElement(msg)
    }
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
            return false
        }
        x.ProcessMessage(stanza)
    }
    
    return true
}