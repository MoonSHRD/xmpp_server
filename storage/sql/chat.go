/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package sql

import (
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ortuman/jackal/model"
	"github.com/ortuman/jackal/xml"
	"github.com/ortuman/jackal/xml/jid"
    "errors"
)

// InsertOrUpdateUser inserts a new user entity into storage,
// or updates it in case it's been previously inserted.
func (s *Storage) InsertOrUpdateChat(c model.Chat) error {
    
    //var columns []string
    //var values []interface{}
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

// FetchUser retrieves from storage a user entity.
func (s *Storage) FetchChat(id int) (*model.Chat, error) {
	q := sq.Select("id", "chatname", "creator", "channel").
		From("chats").
		Where(sq.Eq{"id": id})
	
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
func (s *Storage) FetchChatUsers(id int) (*model.Chat, error) {
	q := sq.Select("username", "admin", "creator", "channel").
		From("chats").
		Where(sq.Eq{"id": id})
	
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
func (s *Storage) DeleteChat(username string) error {
	return s.inTransaction(func(tx *sql.Tx) error {
		var err error
		_, err = sq.Delete("offline_messages").Where(sq.Eq{"username": username}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("roster_items").Where(sq.Eq{"username": username}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("roster_versions").Where(sq.Eq{"username": username}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("private_storage").Where(sq.Eq{"username": username}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("vcards").Where(sq.Eq{"username": username}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		_, err = sq.Delete("users").Where(sq.Eq{"username": username}).RunWith(tx).Exec()
		if err != nil {
			return err
		}
		return nil
	})
}

// UserExists returns whether or not a user exists within storage.
func (s *Storage) UserExists(username string) (bool, error) {
	q := sq.Select("COUNT(*)").From("users").Where(sq.Eq{"username": username})
	var count int
	err := q.RunWith(s.db).QueryRow().Scan(&count)
	switch err {
	case nil:
		return count > 0, nil
	default:
		return false, err
	}
}

// UserExists returns whether or not a user exists within storage.
func (s *Storage) SaveUserNonce(username,nonce string) (error) {
    columns := []string{"nonce", "username", "created_at"}
    values := []interface{}{nonce, username, nowExpr}
    
    var suffix string
    var suffixArgs []interface{}
    suffix = "ON DUPLICATE KEY UPDATE nonce = ?, created_at = NOW()"
    suffixArgs = []interface{}{nonce}
    
    q := sq.Insert("auth_nonce").
        Columns(columns...).
        Values(values...).
        Suffix(suffix, suffixArgs...)
    _, err := q.RunWith(s.db).Exec()
    if err != nil {
        return err
    }
    return nil
}

// UserExists returns whether or not a user exists within storage.
func (s *Storage) LoadUserNonce(username string) (string,error) {
	q := sq.Select("nonce", "created_at").From("auth_nonce").Where(sq.Eq{"username": username})
	
	var nonce string
	var created_at time.Time
    
    err := q.RunWith(s.db).QueryRow().Scan(&nonce,&created_at)
    switch err {
    case nil:
        if nonce=="" {
            return "", errors.New("nonce is empty")
        }
        duration := time.Since(created_at)
        if duration.Minutes()>5 {
            return "", errors.New("too late, retry auth")
        }
        return nonce,nil
    default:
        return "", err
    }
}
