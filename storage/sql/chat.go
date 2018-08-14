/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package sql

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
)

// InsertOrUpdateUser inserts a new user entity into storage,
// or updates it in case it's been previously inserted.
func (s *Storage) InsertOrUpdateChat(c *model.Chat) (int64,error) {
    
    var channel int
    var suffix string
    var suffixArgs []interface{}
    
    if c.Channel {
        channel=1
    } else {
        channel=0
    }
    
    columns := []string{"chatname", "creator", "channel", "created_at", "updated_at"}
    values := []interface{}{c.Chatname, c.Creator, channel, nowExpr, nowExpr}
    
    if c.Id!=0{
        columns=append([]string{"id"},columns...)
        values=append([]interface{}{c.Id},values...)
    
        suffix = "ON DUPLICATE KEY UPDATE chatname = ?, updated_at = NOW()"
        suffixArgs = []interface{}{c.Chatname}
    }
    
    q := sq.Insert("chats").
        Columns(columns...).
        Values(values...).
        Suffix(suffix, suffixArgs...)
    res, err := q.RunWith(s.db).Exec()
    id,err:=res.LastInsertId()
    return id,err
}

func (s *Storage) InsertChatUser(chat_id int64,username string,admin bool) error {
    
    //var columns []string
    //var values []interface{}
    var i_admin int
    
    if admin {
        i_admin=1
    } else {
        i_admin=0
    }
    
    columns := []string{"chat_id", "username", "admin", "created_at"}
    values := []interface{}{chat_id, username, i_admin, nowExpr}
    
    //var suffix string
    //var suffixArgs []interface{}
    
    //suffix = "ON DUPLICATE KEY IGNORE"
    
    q := sq.Insert("chats_users").
        Columns(columns...).
        Values(values...)
    _, err := q.RunWith(s.db).Exec()
    //fmt.Println(id)
    return err
}

func (s *Storage) DeleteChatUser(chat_id int64,username string) error {
    return s.inTransaction(func(tx *sql.Tx) error {
        var err error
        _, err = sq.Delete("chats_users").Where(sq.Eq{"chat_id": chat_id,"username":username}).RunWith(tx).Exec()
        if err != nil {
            return err
        }
        _, err = sq.Delete("chats_msgs").Where(sq.Eq{"chat_id": chat_id,"username":username}).RunWith(tx).Exec()
        if err != nil {
            return err
        }
        return nil
    })
}

// FetchUser retrieves from storage a user entity.
func (s *Storage) FetchChat(chat_id int64) (*model.Chat, error) {
	q := sq.Select("id", "chatname", "creator", "channel").
		From("chats").
		Where(sq.Eq{"id": chat_id})
	
	var c model.Chat

	err := q.RunWith(s.db).QueryRow().Scan(&c.Id, &c.Chatname, &c.Creator, &c.Channel)
	switch err {
	case nil:
		return &c, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// FetchUser retrieves from storage a user entity.
func (s *Storage) FetchChatUsers(chat_id int64) (model.ChatUsers, error) {
	q := sq.Select("username", "admin").
		From("chats_users").
		Where(sq.Eq{"chat_id": chat_id})

	rows,err := q.RunWith(s.db).Query()
	switch err {
	case nil:
	    
	    users := model.ChatUsers{}
	    var username string
	    var admin int
	    for rows.Next() {
            rows.Scan(&username,&admin)
            users[username]=model.ChatUser{username,admin}
        }
		return users, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

// DeleteUser deletes a user entity from storage.
func (s *Storage) DeleteChat(chat_id int64) error {
	return s.inTransaction(func(tx *sql.Tx) error {
		var err error
		_, err = sq.Delete("chats_msg").Where(sq.Eq{"chat_id": chat_id}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("chats_users").Where(sq.Eq{"chat_id": chat_id}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("chats").Where(sq.Eq{"id": chat_id}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		return nil
	})
}

// UserExists returns whether or not a user exists within storage.
func (s *Storage) ChatExists(chat_id int64) (bool, error) {
	q := sq.Select("COUNT(*)").From("chats").Where(sq.Eq{"id": chat_id})
	var count int
	err := q.RunWith(s.db).QueryRow().Scan(&count)
	switch err {
	case nil:
		return count > 0, nil
	default:
		return false, err
	}
}