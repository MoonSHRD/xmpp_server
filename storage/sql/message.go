package sql

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
	"strings"
)

func (s *Storage) WriteMsgToDB(chat_id, sender, msg string, isOnline, isFile int) (int64, string, error) {
	q := sq.Insert("messages").
		Columns("chat_id", "sender", "message", "created_at", "updated_at", "delivered", "files").
		Values(chat_id, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"), isOnline, isFile)
		res, err := q.RunWith(s.db).Exec()
		if err != nil {
			return 0, "", err
		}
		id, _ := res.LastInsertId()
		date_query := sq.Select("created_at").From("messages").Where("id = ?", id)
		res_date, _ := date_query.RunWith(s.db).Query()
		var _date string

		for res_date.Next() {
				res_date.Scan(&_date)
			}
		_date = strings.Replace(_date, "T", " ", -1)
		return id, _date, nil
}

func (s *Storage) WriteFileToDB(file model.File, msg_id int64) error{
	q := sq.Insert("files").
		Columns("message_id", "hash", "type", "name").
		Values(msg_id, file.Hash, file.Type, file.Name)
	_, err := q.RunWith(s.db).Exec()

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetMsgsFromDB(chat_id string) ([]model.Message, error) {
	var list_messages []model.Message
	q := sq.Select("id", "sender", "message", "created_at", "files").From("messages").Where("chat_id = ?", chat_id).OrderBy("created_at")
	records, err:= q.RunWith(s.db).Query()
	if err!= nil {
		return nil, err
	}
	for records.Next() {
		message := model.Message{}
		records.Scan(&message.Id, &message.Sender, &message.Message, &message.Time, &message.File)
		list_messages = append(list_messages, message)
	}
	return list_messages, nil

}

func (s *Storage) GetFilesFromDB(msg_id int64) ([]model.File, error) {
	var list_files []model.File
	q := sq.Select("hash", "type", "name").From("files").Where("message_id = ?", msg_id)
	records, err:= q.RunWith(s.db).Query()
	if err!= nil {
		return nil, err
	}
	for records.Next() {
		file := model.File{}
		records.Scan(&file.Hash, &file.Type, &file.Name)
		list_files = append(list_files, file)
	}
	return list_files, nil
}
