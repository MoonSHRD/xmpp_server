package model

import (
	"time"
)

type Message struct {
	Id int64
	ChatId string
	Sender string
	Time time.Time
	Message string
	File int
}
