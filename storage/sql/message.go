package sql

import (
	sq "github.com/Masterminds/squirrel"
)


//var (
//	nowExpr = sq.Expr("NOW()")
//)

func (s *Storage) write_msg_to_db(recipient, sender, msg string) error {
	q := sq.Insert("jackal").
		Columns("recipient", "sender", "message", "created_at", "updated_at").
		Values(recipient, sender, msg, sq.Expr("NOW()"), sq.Expr("NOW()"))
	_, err := q.RunWith(s.db).Exec()
    if err != nil {
		return err
	}
	return nil
}