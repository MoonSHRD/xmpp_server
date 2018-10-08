package sql

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
	"strings"
)

func (s *Storage) WriteMsgToDB(chat_id, sender, msg string, isOnline int) (int64, string, error) {
	//check, _ := (*Storage).UserExists(s, chat_id)
	//if check == false {
	//	return 0, nil
	//}
	q := sq.Insert("messages").
		Columns("chat_id", "sender", "message", "created_at", "updated_at", "delivered").
		Values(chat_id, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"), isOnline)
		res, err := q.RunWith(s.db).Exec()
		if err != nil {
			return 0, "", err
		}
		id, _ := res.LastInsertId()
		//date:= sq.Expr("NOW()")
		date_query := sq.Select("created_at").From("messages").Where("id = ?", id)
		res_date, _ := date_query.RunWith(s.db).Query()
		var _date string

		for res_date.Next() {
				res_date.Scan(&_date)
			}
		_date = strings.Replace(_date, "T", " ", -1)
		return id, _date, nil
}

func (s *Storage) GetMsgFromDB(chat_id string) ([]model.Message, error) {
	var list_messages []model.Message
	q := sq.Select("sender", "message", "created_at").From("messages").Where("chat_id = ?", chat_id).OrderBy("created_at")
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

//func (s *Storage) GetDateOfMsg(chat_id string) (string) {
//	date_query := sq.Select("created_at").From("messages").Where("id = ?", id)
//	res_date, _ := date_query.RunWith(s.db).Query()
//	var _date string
//
//	for res_date.Next() {
//		res_date.Scan(&_date)
//	}
//}
