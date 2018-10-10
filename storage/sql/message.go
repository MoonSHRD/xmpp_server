package sql

import (
	sq "github.com/Masterminds/squirrel"
)


//var (
//	nowExpr = sq.Expr("NOW()")
//)

func (s *Storage) WriteMsgToDB(recipient, sender, msg string) (bool, error) {
	check, _ := (*Storage).UserExists(s, recipient)
	if check == false {
		print("Non exist")
		return false, nil
	}
	q := sq.Insert("users_messages").
		Columns("recipient", "sender", "message", "created_at", "updated_at").
		Values(recipient, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"))
		_, err := q.RunWith(s.db).Exec()
		if err != nil {
			return true, err
		}
		return true, nil
}