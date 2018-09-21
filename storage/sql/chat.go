/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package sql

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
    "fmt"
)

// InsertOrUpdateUser inserts a new user entity into storage,
// or updates it in case it's been previously inserted.
func (s *Storage) InsertOrUpdateChat(c *model.Chat) (string, error) {
    
    var channel int
    var suffix string
    var suffixArgs []interface{}
    
    if c.Channel {
        channel=1
    } else {
        channel=0
    }
    
    columns := []string{"id", "chatname", "creator", "channel", "created_at", "updated_at", "avatar"}
    values := []interface{}{c.Id, c.Chatname, c.Creator, channel, nowExpr, nowExpr, c.Avatar}
    
    if c.Id!= ""{
        //columns=append([]string{"id"},columns...)
        //values=append([]interface{}{c.Id},values...)
    
        suffix = "ON DUPLICATE KEY UPDATE chatname = ?, updated_at = NOW(), avatar = ?"
        suffixArgs = []interface{}{c.Chatname,c.Avatar}
    }
    
    q := sq.Insert("chats").
        Columns(columns...).
        Values(values...).
        Suffix(suffix, suffixArgs...)
    _, err := q.RunWith(s.db).Exec()
    if err!=nil {
        fmt.Println(err)
    }
    //id,err:=res.LastInsertId()
    //if err!=nil {
    //    fmt.Println(err)
    //}
    return c.Id, err
}

func (s *Storage) InsertChatUser(chat_id string,username string,admin bool) error {
    
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

func (s *Storage) DeleteChatUser(chat_id string,username string) error {
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
func (s *Storage) FetchChat(chat_id string) (*model.Chat, error) {
	q := sq.Select("id", "chatname", "creator", "channel", "avatar").
		From("chats").
		Where(sq.Eq{"id": chat_id})
	
	var c model.Chat

	err := q.RunWith(s.db).QueryRow().Scan(&c.Id, &c.Chatname, &c.Creator, &c.Channel, &c.Avatar)
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
func (s *Storage) FetchChatUsers(chat_id string) (model.ChatUsers, error) {
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
func (s *Storage) DeleteChat(chat_id string) error {
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
func (s *Storage) ChatExists(chat_name string) (bool, error) {
	q := sq.Select("COUNT(*)").From("chats").Where(sq.Eq{"chatname": chat_name})
	var count int
	err := q.RunWith(s.db).QueryRow().Scan(&count)
	switch err {
	case nil:
		return count > 0, nil
	default:
		return false, err
	}
}

func (s *Storage) FindGroups(chat_name string) []model.Chat{
	q := sq.Select("id", "chatname", "creator", "channel", "avatar").From("chats").Where("chatname LIKE ? or id = ?", "%" + chat_name + "%", chat_name)
	records, _:= q.RunWith(s.db).Query()
	var list_chats []model.Chat
	for records.Next(){
	    chat:=model.Chat{}
		records.Scan(&chat.Id, &chat.Chatname, &chat.Creator, &chat.Channel, &chat.Avatar)
		list_chats = append(list_chats, chat)
	}
	return list_chats
}
