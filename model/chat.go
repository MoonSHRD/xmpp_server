package model

type Chat struct {
    Id string
    Chatname string
    Creator string
    Channel bool
    Avatar string
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
