package models

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/asaskevich/govalidator"
	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/tommy351/maji.moe/errors"
)

type User struct {
	ID                 int64     `db:"id" json:"id"`
	Name               string    `db:"name" json:"name"`
	Password           string    `db:"password" json:"-"`
	Email              string    `db:"email" json:"email"`
	Avatar             string    `db:"avatar" json:"avatar"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
	Activated          bool      `db:"activated" json:"activated"`
	ActivationToken    string    `db:"activation_token" json:"-"`
	PasswordResetToken string    `db:"password_reset_token" json:"-"`
	FacebookID         string    `db:"facebook_id" json:"-"`
	TwitterID          string    `db:"twitter_id" json:"-"`
	GoogleID           string    `db:"google_id" json:"-"`
	GithubID           string    `db:"github_id" json:"-"`
}

func (data *User) Validate() error {
	if govalidator.IsNull(data.Name) {
		return errors.New("name", errors.Required, "Name is required")
	}

	if !govalidator.IsEmail(data.Email) {
		return errors.New("email", errors.Email, "Email is invalid")
	}

	data.Name = govalidator.Trim(data.Name, "")

	return nil
}

func (data *User) PreInsert(s gorp.SqlExecutor) error {
	if err := data.Validate(); err != nil {
		return err
	}

	// Check whether email has been taken
	if count, _ := s.SelectInt("SELECT count(id) FROM users WHERE email=?", data.Email); count > 0 {
		return errors.New("email", errors.EmailUsed, "Email has been taken")
	}

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *User) PreUpdate(s gorp.SqlExecutor) error {
	if err := data.Validate(); err != nil {
		return err
	}

	data.UpdatedAt = time.Now()
	return nil
}

func (data *User) PreDelete(s gorp.SqlExecutor) error {
	if _, err := s.Exec("DELETE FROM domains WHERE user_id=?", data.ID); err != nil {
		return err
	}

	if _, err := s.Exec("DELETE FROM tokens WHERE user_id=?", data.ID); err != nil {
		return err
	}

	return nil
}

func (data *User) Gravatar() {
	h := md5.New()
	io.WriteString(h, data.Email)

	data.Avatar = "//www.gravatar.com/avatar/" + hex.EncodeToString(h.Sum(nil))
}

func (data *User) GeneratePassword(password string) error {
	if govalidator.IsNull(password) {
		return errors.New("password", errors.Required, "Password is required")
	}

	if !govalidator.IsByteLength(password, 6, 50) {
		return errors.New("password", errors.Length, "The length of password must be between 6-50")
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		return err
	} else {
		data.Password = string(hash)
	}

	return nil
}

func (data *User) Authenticate(password string) error {
	if govalidator.IsByteLength(password, 6, 50) {
		return errors.New("password", errors.Length, "The length of password must be between 6-50")
	}

	if govalidator.IsNull(data.Password) {
		return errors.New("password", errors.PasswordUnset, "Password has not been set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password)); err != nil {
		return errors.API{
			Status:  http.StatusUnauthorized,
			Field:   "password",
			Code:    errors.WrongPassword,
			Message: "Password is wrong",
		}
	}

	return nil
}

func (data *User) SetActivated(activated bool) {
	if activated {
		data.Activated = true
	} else {
		data.Activated = false
		data.ActivationToken = uniuri.NewLen(32)
	}
}

func (data *User) SendActivationMail() {
	//
}
