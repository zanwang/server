package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

// Record model
type Record struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Type        string `db:"type" json:"type"`
	Destination string `db:"destination" json:"destination"`
	CreatedAt   int64  `db:"created_at" json:"created_at"`
	UpdatedAt   int64  `db:"updated_at" json:"updated_at"`
	DomainID    int64  `db:"domain_id" json:"domain_id"`
	UserID      int64  `db:"user_id" json:"user_id"`
}

func (data *Record) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now().UnixNano()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *Record) PreUpdate(s gorp.SqlExecutor) error {
	data.UpdatedAt = time.Now().UnixNano()
	return nil
}
