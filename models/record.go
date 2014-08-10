package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

var RecordType = []string{"A", "CNAME", "MX", "TXT", "SPF", "AAAA", "NS", "LOC"}

// Record model
type Record struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Type      string `db:"type" json:"type"`
	Content   string `db:"content" json:"content"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	DomainID  int64  `db:"domain_id" json:"domain_id"`
	TTL       uint   `db:"ttl" json:"ttl"`
	Priority  uint   `db:"priority" json:"priority"`
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
