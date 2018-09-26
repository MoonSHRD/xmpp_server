package sql

import (
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) WriteMsgToDB(recipient, sender, msg string, isOnline int) (int64, error) {
	check, _ := (*Storage).UserExists(s, recipient)
	if check == false {
		return 0, nil
	}
	q := sq.Insert("messages").
		Columns("chat_id", "sender", "message", "created_at", "updated_at", "delivered").
		Values(recipient, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"), isOnline)
		res, err := q.RunWith(s.db).Exec()
		if err != nil {
			return 0, err
		}
		id, _ := res.LastInsertId()
		return id, nil
}
