package xep0045

import (
    "github.com/ortuman/jackal/xml"
    "github.com/ortuman/jackal/stream"
    "github.com/ortuman/jackal/module/xep0030"
    "github.com/ortuman/jackal/storage"
    "github.com/ortuman/jackal/model"
    "strconv"
)

const chatNamespace = "http://jabber.org/protocol/muc"

type RegisterChat struct {
    stm        stream.C2S
}

func (x *RegisterChat) RegisterDisco(discoInfo *xep0030.DiscoInfo) {
    // register disco feature
    discoInfo.Entity(x.stm.Domain(), "").AddFeature(chatNamespace)
}

func New(stm stream.C2S) *RegisterChat {
    return &RegisterChat{
        stm: stm,
    }
}


func (x *RegisterChat) CreateChat(presence *xml.Presence) {
    to:=presence.ToJID()
    from:=presence.FromJID().String()
    chat:=model.Chat{Chatname:to.Node(),Channel:false,Creator:from}
    var err error
    chat.Id, err = storage.Instance().InsertOrUpdateChat(&chat)
    storage.Instance().InsertChatUser(chat.Id,presence.FromJID().Node(),true)
    if err != nil {
        x.stm.SendElement(presence.NotAllowedError())
    } else {
        item_elem:=xml.NewElementName("item")
        item_elem.SetAttribute("affiliation","owner")
        item_elem.SetAttribute("role","moderator")
        
        x_elem:=xml.NewElementName("x")
        x_elem.SetNamespace(chatNamespace+"#user")
        x_elem.AppendElement(item_elem)
        
        elem:=xml.NewElementName("presence")
        elem.SetFrom(strconv.Itoa(int(chat.Id))+"@localhost/"+to.Node())
        elem.SetTo(from)
        elem.AppendElement(x_elem)
        
        x.stm.SendElement(elem)
    }
}

func (x *RegisterChat) ProcessMessage(msg *xml.Message) {
    
    node := msg.ToJID().Node()
    if node == "" {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    to,err:=strconv.ParseInt(node,10,64)
    if err != nil {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    exist,err := storage.Instance().ChatExists(to)
    
    if err != nil {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    if !exist {
        x.stm.SendElement(msg.BadRequestError())
        return
    }
    
    chat_u,_ := storage.Instance().FetchChatUsers(to)
    //chat,_ := storage.Instance().FetchChat(to)
    
    msg.SetFrom(msg.To())
    for _,username := range chat_u {
        msg.SetTo(username)
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
        x.CreateChat(stanza)

    case *xml.Message:
    
        if !stanza.IsGroupChat() {
            return false
        }
        x.ProcessMessage(stanza)
    }
    
    return true
}