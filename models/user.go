package models

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/asaskevich/govalidator"
	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/majimoe/server/config"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/util"
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

const (
	mailSender    = "maji.moe <noreply@maji.moe>"
	mailTitle     = "Activate your account"
	mailRecipient = "%s <%s>"
)

var mailTmpl *template.Template

func init() {
	baseDir := config.BaseDir
	mailTmpl = template.Must(template.ParseFiles(
		path.Join(baseDir, "views", "email", "activation.html"),
	))
}

func (data *User) Validate(s gorp.SqlExecutor) error {
	if govalidator.IsNull(data.Name) {
		return errors.New("name", errors.Required, "Name is required")
	}

	if !govalidator.IsEmail(data.Email) {
		return errors.New("email", errors.Email, "Email is invalid")
	}

	var user User

	if err := s.SelectOne(&user, "SELECT id FROM users WHERE email=?", data.Email); err == nil {
		if data.ID != user.ID {
			return errors.New("email", errors.EmailUsed, "Email has been taken")
		}
	}

	data.Name = govalidator.Trim(data.Name, "")

	return nil
}

func (data *User) PreInsert(s gorp.SqlExecutor) error {
	if err := data.Validate(s); err != nil {
		return err
	}

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return nil
}

func (data *User) PreUpdate(s gorp.SqlExecutor) error {
	if err := data.Validate(s); err != nil {
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
	if !govalidator.IsByteLength(password, 6, 50) {
		return errors.New("password", errors.Length, "The length of password must be between 6-50")
	}

	if govalidator.IsNull(data.Password) {
		return errors.API{
			Status:  http.StatusUnauthorized,
			Field:   "password",
			Code:    errors.PasswordUnset,
			Message: "Password has not been set",
		}
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
	if data.Activated || !config.Config.EmailActivation {
		return
	}

	var buf bytes.Buffer
	err := mailTmpl.Execute(&buf, map[string]interface{}{
		"User": data,
	})

	if err != nil {
		log.Println(err)
		return
	}

	recipient := fmt.Sprintf(mailRecipient, data.Name, data.Email)
	msg := util.Mailgun.NewMessage(mailSender, mailTitle, buf.String(), recipient)

	if _, _, err := util.Mailgun.Send(msg); err != nil {
		log.Println(err)
	}
}
