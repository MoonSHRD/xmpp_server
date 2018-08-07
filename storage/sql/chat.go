/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package sql

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
	"time"
)

// InsertOrUpdateUser inserts a new user entity into storage,
// or updates it in case it's been previously inserted.
func (s *Storage) InsertOrUpdateChat(c model.Chat) error {
    
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
    
    if c.Id!=""{
        columns=append([]string{"id"},columns...)
        values=append([]interface{}{c.Id},values...)
    
        suffix = "ON DUPLICATE KEY UPDATE chatname = ?, updated_at = NOW()"
        suffixArgs = []interface{}{c.Chatname}
    }
    
    q := sq.Insert("users").
        Columns(columns...).
        Values(values...).
        Suffix(suffix, suffixArgs...)
    _, err := q.RunWith(s.db).Exec()
    return err
}

func (s *Storage) InsertChatUser(chat_id int,username string,admin bool) error {
    
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
    
    var suffix string
    var suffixArgs []interface{}
    
    suffix = "ON DUPLICATE KEY IGNORE"
    
    q := sq.Insert("users").
        Columns(columns...).
        Values(values...).
        Suffix(suffix, suffixArgs...)
    _, err := q.RunWith(s.db).Exec()
    return err
}

func (s *Storage) DeleteChatUser(chat_id int,username string) error {
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
func (s *Storage) FetchChat(chat_id int) (*model.Chat, error) {
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
func (s *Storage) FetchChatUsers(chat_id int) (*model.Chat, error) {
	q := sq.Select("username", "admin").
		From("chats").
		Where(sq.Eq{"chat_id": chat_id})
	
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

// DeleteUser deletes a user entity from storage.
func (s *Storage) DeleteChat(chat_id int) error {
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
		_, err = sq.Delete("chats").Where(sq.Eq{"chat_id": chat_id}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		return nil
	})
}

// UserExists returns whether or not a user exists within storage.
func (s *Storage) ChatExists(chat_id string) (bool, error) {
	q := sq.Select("COUNT(*)").From("chats").Where(sq.Eq{"chat_id": chat_id})
	var count int
	err := q.RunWith(s.db).QueryRow().Scan(&count)
	switch err {
	case nil:
		return count > 0, nil
	default:
		return false, err
	}
}

func (s *Storage) AddMessage(chat_id int,username, message string) {
	db, err := sql.Open("mysql", "jackal:password@/jackal")  // Change Password!!!!!

		if err != nil {
			panic(err)
		}
		defer db.Close()
		created_at := time.Now()
		updated_at := time.Now()

		result, err := db.Exec("insert into jackal.chats_msgs (chat_id, username, msg, created_at, updated_at) values (?, ?)",
			 chat_id, username, message, created_at, updated_at)
		if err != nil{
			panic(err)
		}
		print(result)
}

func (s *Storage) EditMessage(chat_id int,username, old_msg, new_msg string) {
	db, err := sql.Open("mysql", "jackal:password@/jackal")  // Change Password!!!!!

		if err != nil {
			panic(err)
		}
		defer db.Close()
		updated_at := time.Now()

		result, err := db.Exec("update chats_msg message = ?, updated_at = ? where chat_id = ? and username = ? and msg = ?",
			new_msg, updated_at, chat_id, username, old_msg)
		if err != nil{
			panic(err)
		}
		print(result)
}

func (s *Storage) DelMessage(chat_id int,username, msg string) {
	db, err := sql.Open("mysql", "jackal:password@/jackal")  // Change Password!!!!!

		if err != nil {
			panic(err)
		}
		defer db.Close()

		result, err := db.Exec("delete from chats_msg where chat_id=? and username=? and msg=?",
			chat_id, username, msg)
		if err != nil{
			panic(err)
		}
		print(result)
}