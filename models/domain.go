package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

// Domain model
type Domain struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Public    bool      `db:"public" json:"public"`
}

func (data *Domain) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *Domain) PreUpdate(s gorp.SqlExecutor) error {
	data.UpdatedAt = time.Now()
	return nil
}

func (data *Domain) PreDelete(s gorp.SqlExecutor) error {
	if _, err := s.Exec("DELETE FROM records WHERE domain_id=?", data.ID); err != nil {
		return err
	}

	return nil
}
