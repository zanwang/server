package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

var RecordType = []string{"A", "CNAME", "MX", "TXT", "SPF", "AAAA", "NS", "LOC"}

// Record model
type Record struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DomainID  int64     `db:"domain_id" json:"domain_id"`
	TTL       uint      `db:"ttl" json:"ttl"`
	Priority  uint      `db:"priority" json:"priority"`
}

func (data *Record) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *Record) PreUpdate(s gorp.SqlExecutor) error {
	data.UpdatedAt = time.Now()
	return nil
}
