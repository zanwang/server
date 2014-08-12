package models

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"time"

	"github.com/coopernurse/gorp"
)

// User model
type User struct {
	ID              int64  `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
	Password        string `db:"password" json:"-"`
	Email           string `db:"email" json:"email"`
	DisplayName     string `db:"display_name" json:"display_name"`
	Avatar          string `db:"avatar" json:"avatar"`
	CreatedAt       int64  `db:"created_at" json:"created_at"`
	UpdatedAt       int64  `db:"updated_at" json:"updated_at"`
	Activated       bool   `db:"activated" json:"activated"`
	ActivationToken string `db:"activation_token" json:"-"`
	LoggedIn        bool   `db:"-" json:"-"`
}

func (data *User) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now().UnixNano()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *User) PreUpdate(s gorp.SqlExecutor) error {
	data.UpdatedAt = time.Now().UnixNano()
	return nil
}

func (data *User) PreDelete(s gorp.SqlExecutor) error {
	if _, err := s.Exec("DELETE FROM domains WHERE user_id=?", data.ID); err != nil {
		return err
	}

	return nil
}

func (data *User) Gravatar() {
	h := md5.New()
	io.WriteString(h, data.Email)

	data.Avatar = "http://www.gravatar.com/avatar/" + hex.EncodeToString(h.Sum(nil))
}

func (data *User) SendActivationEmail() {
	//
}
