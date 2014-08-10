package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

// Domain model
type Domain struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	UserID    int64  `db:"user_id" json:"user_id"`
	Public    bool   `db:"public" json:"public"`
}

func (data *Domain) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now().UnixNano()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *Domain) PreUpdate(s gorp.SqlExecutor) error {
	data.UpdatedAt = time.Now().UnixNano()
	return nil
}
