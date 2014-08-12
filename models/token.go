package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

// Token model
type Token struct {
	ID         int64     `db:"id" json:"-"`
	UserID     int64     `db:"user_id" json:"-"`
	Key        string    `db:"key" json:"key"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	Authorized bool      `db:"-" json:"-"`
}

func (data *Token) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *Token) PreUpdate(s gorp.SqlExecutor) error {
	data.UpdatedAt = time.Now()
	return nil
}
