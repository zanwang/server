package models

import (
	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
)

const (
	tokenExpiry = 60 * 60 * 24 * 7
)

// Token model
type Token struct {
	ID        int64  `db:"id" json:"-"`
	UserID    int64  `db:"user_id" json:"user_id"`
	Key       string `db:"key" json:"key"`
	CreatedAt int64  `db:"created_at" json:"-"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	ExpiredAt int64  `db:"expired_at" json:"expired_at"`
}

func (data *Token) PreInsert(s gorp.SqlExecutor) error {
	now := Now()
	data.Key = uniuri.NewLen(32)
	data.CreatedAt = now
	data.UpdatedAt = now
	data.ExpiredAt = now + tokenExpiry
	return nil
}

func (data *Token) PreUpdate(s gorp.SqlExecutor) error {
	now := Now()
	data.UpdatedAt = now
	data.ExpiredAt = now + tokenExpiry
	return nil
}
