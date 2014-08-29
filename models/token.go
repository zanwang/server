package models

import (
	"time"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
)

const (
	tokenExpiry = time.Hour * 24 * 7 // 7 days
)

// Token model
type Token struct {
	ID        int64     `db:"id" json:"-"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Key       string    `db:"key" json:"key"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
}

func (data *Token) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	data.Key = uniuri.NewLen(32)
	data.CreatedAt = now
	data.UpdatedAt = now
	data.ExpiredAt = now.Add(tokenExpiry)
	return nil
}

func (data *Token) PreUpdate(s gorp.SqlExecutor) error {
	now := time.Now()
	data.UpdatedAt = now
	data.ExpiredAt = now.Add(tokenExpiry)
	return nil
}
