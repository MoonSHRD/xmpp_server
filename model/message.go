package model

import (
	"time"
)

type Message struct {
	ChatId string
	Sender string
	Time time.Time
	Message string
}
