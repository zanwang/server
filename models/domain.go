package models

type Domain struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	UserID    int64  `db:"user_id"`
	Public    bool   `db:"public"`
}
