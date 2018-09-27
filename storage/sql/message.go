package sql

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
)

func (s *Storage) WriteMsgToDB(chat_id, sender, msg string, isOnline int) (int64, error) {
	//check, _ := (*Storage).UserExists(s, chat_id)
	//if check == false {
	//	return 0, nil
	//}
	q := sq.Insert("messages").
		Columns("chat_id", "sender", "message", "created_at", "updated_at", "delivered").
		Values(chat_id, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"), isOnline)
		res, err := q.RunWith(s.db).Exec()
		if err != nil {
			return 0, err
		}
		id, _ := res.LastInsertId()
		return id, nil
}

func (s *Storage) GetMsgFromDB(chat_id string) ([]model.Message, error) {
	var list_messages []model.Message
	q := sq.Select("sender", "message", "created_at").From("messages").Where("chat_id = ?", chat_id).OrderBy("message")
	records, err:= q.RunWith(s.db).Query()
	if err!= nil {
		return nil, err
	}
	for records.Next() {
		message := model.Message{}
		records.Scan(&message.Sender, &message.Message, &message.Time)
		list_messages = append(list_messages, message)
	}
	return list_messages, nil

}
