package model

type Chat struct {
    Id string
    Chatname string
    Creator string
    Type string
    Avatar string
}

type ChatUser struct {
    Username string
    Role string
}

type ChatUsers map[string]ChatUser

type File struct {
    Hash string
    Type string
    Name string
}

//func (chat Chat) IsChannel() string {
//    channel:="0"
//    if chat.Channel {
//        channel="1"
//    }
//    return channel
//}
