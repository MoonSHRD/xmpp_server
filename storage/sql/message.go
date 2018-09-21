package sql

import (
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) WriteMsgToDB(recipient, sender, msg string) (int64, error) {
	check, _ := (*Storage).UserExists(s, recipient)
	if check == false {
		print("Non exist")
		return 0, nil
	}
	q := sq.Insert("users_messages").
		Columns("recipient", "sender", "message", "created_at", "updated_at").
		Values(recipient, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"))
		res, err := q.RunWith(s.db).Exec()
		if err != nil {
			return 0, err
		}
		id, _ := res.LastInsertId()
		return id, nil
}
