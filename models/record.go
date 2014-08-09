package models

type Record struct {
	ID          int64  `db:"id"`
	Type        string `db:"type"`
	Subdomain   string `db:"subdomain"`
	Destination string `db:"destination"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
	DomainID    int64  `db:"domain_id"`
}
