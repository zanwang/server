package models

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/dchest/uniuri"
	"github.com/majimoe/server/config"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/util"
)

const (
	mailSender             = "maji.moe <noreply@maji.moe>"
	mailRecipient          = "%s <%s>"
	activationMailTitle    = "Activate your account"
	passwordResetMailTitle = "Reset your password"
)

var (
	activationMailTmpl, passwordResetMailTmpl *template.Template
)

func init() {
	baseDir := config.BaseDir
	emailViewDir := path.Join(baseDir, "views", "email")

	activationMailTmpl = template.Must(template.ParseFiles(
		path.Join(emailViewDir, "activation.html"),
	))
	passwordResetMailTmpl = template.Must(template.ParseFiles(
		path.Join(emailViewDir, "password_reset.html"),
	))
}

type User struct {
	Id                 int64
	Name               string
	Password           string
	Email              string
	Avatar             string
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Activated          bool
	ActivationToken    string
	PasswordResetToken string
	FacebookId         string
	GoogleId           string
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         u.Id,
		"name":       u.Name,
		"email":      u.Email,
		"avatar":     u.Avatar,
		"created_at": ISOTime(u.CreatedAt),
		"updated_at": ISOTime(u.UpdatedAt),
		"activated":  u.Activated,
	})
}

func (u *User) PublicProfile() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.Id,
		"name":       u.Name,
		"avatar":     u.Avatar,
		"created_at": ISOTime(u.CreatedAt),
		"updated_at": ISOTime(u.UpdatedAt),
	}
}

func (u *User) BeforeSave() error {
	if govalidator.IsNull(u.Name) {
		return errors.New("name", errors.Required, "Name is required")
	}

	if !govalidator.IsEmail(u.Email) {
		return errors.New("email", errors.Email, "Email is invalid")
	}

	if len(u.Name) > 100 {
		return errors.New("name", errors.MaxLength, "Maximum length of name is 100")
	}

	u.Name = govalidator.Trim(u.Name, "")
	u.UpdatedAt = time.Now().UTC()

	return nil
}

func (u *User) BeforeCreate() error {
	u.CreatedAt = time.Now().UTC()
	return nil
}

func (u *User) Gravatar() {
	h := md5.New()
	io.WriteString(h, u.Email)

	u.Avatar = "//www.gravatar.com/avatar/" + hex.EncodeToString(h.Sum(nil))
}

func (u *User) GeneratePassword(password string) error {
	if govalidator.IsNull(password) {
		return errors.New("password", errors.Required, "Password is required")
	}

	if len(password) < 6 {
		return errors.New("password", errors.MinLength, "Minimum length of password is 6")
	}

	if len(password) > 50 {
		return errors.New("password", errors.MaxLength, "Maximum length of password is 50")
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		return err
	} else {
		u.Password = string(hash)
	}

	return nil
}

func (u *User) Authenticate(password string) error {
	if len(password) < 6 {
		return errors.New("password", errors.MinLength, "Minimum length of password is 6")
	}

	if len(password) > 50 {
		return errors.New("password", errors.MaxLength, "Maximum length of password is 50")
	}

	if govalidator.IsNull(u.Password) {
		return &errors.API{
			Status:  http.StatusUnauthorized,
			Field:   "password",
			Code:    errors.PasswordUnset,
			Message: "Password has not been set",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return &errors.API{
			Status:  http.StatusUnauthorized,
			Field:   "password",
			Code:    errors.WrongPassword,
			Message: "Password is wrong",
		}
	}

	return nil
}

func (u *User) SetActivated(activated bool) {
	if activated {
		u.Activated = true
	} else {
		u.Activated = false
		u.ActivationToken = uniuri.NewLen(32)
	}
}

func (u *User) SendActivationMail() {
	if u.Activated || !config.Config.EmailActivation {
		return
	}

	var buf bytes.Buffer

	if err := activationMailTmpl.Execute(&buf, map[string]interface{}{
		"User": u,
	}); err != nil {
		log.Println(err)
		return
	}

	recipient := fmt.Sprintf(mailRecipient, u.Name, u.Email)
	msg := util.Mailgun.NewMessage(mailSender, activationMailTitle, buf.String(), recipient)

	if _, _, err := util.Mailgun.Send(msg); err != nil {
		log.Println(err)
	}
}

func (u *User) SendPasswordResetMail() {
	if !config.Config.EmailActivation {
		return
	}

	var buf bytes.Buffer

	if err := passwordResetMailTmpl.Execute(&buf, map[string]interface{}{
		"User": u,
	}); err != nil {
		log.Println(err)
		return
	}

	recipient := fmt.Sprintf(mailRecipient, u.Name, u.Email)
	msg := util.Mailgun.NewMessage(mailSender, passwordResetMailTitle, buf.String(), recipient)

	if _, _, err := util.Mailgun.Send(msg); err != nil {
		log.Println(err)
	}
}
