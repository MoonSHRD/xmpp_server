package model

type Chat struct {
    Id int64
    Chatname string
    Creator string
    Channel bool
}

type ChatUser struct {
    Username string
    Admin int
}

type ChatUsers map[string]ChatUser

func (chat Chat) IsChannel() string {
    channel:="0"
    if chat.Channel {
        channel="1"
    }
    return channel
}