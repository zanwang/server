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

	if !govalidator.IsByteLength(password, 6, 50) {
		return errors.New("password", errors.Length, "The length of password must be between 6-50")
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		return err
	} else {
		u.Password = string(hash)
	}

	return nil
}

func (u *User) Authenticate(password string) error {
	if !govalidator.IsByteLength(password, 6, 50) {
		return errors.New("password", errors.Length, "The length of password must be between 6-50")
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
	err := mailTmpl.Execute(&buf, map[string]interface{}{
		"User": u,
	})

	if err != nil {
		log.Println(err)
		return
	}

	recipient := fmt.Sprintf(mailRecipient, u.Name, u.Email)
	msg := util.Mailgun.NewMessage(mailSender, mailTitle, buf.String(), recipient)

	if _, _, err := util.Mailgun.Send(msg); err != nil {
		log.Println(err)
	}
}
