package models

type User struct {
	ID              int64  `db:"id"`
	Username        string `db:"username"`
	Password        string `db:"password"`
	Email           string `db:"email"`
	CreatedAt       int64  `db:"created_at"`
	UpdatedAt       int64  `db:"updated_at"`
	Activated       bool   `db:"activated"`
	DisplayName     string `db:"display_name"`
	ActivationToken string `db:"activation_token"`
	LoggedIn        bool   `db:"-"`
}
