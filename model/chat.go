package model

type Chat struct {
    Id string
    Chatname string
    Creator string
    Channel bool
}

func NewChat(name,username string,channel bool) Chat {
    chat:=Chat{Chatname:name,Creator:username,Channel:channel}
    return chat
}
