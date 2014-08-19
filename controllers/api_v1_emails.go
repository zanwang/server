package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/coopernurse/gorp"
	"github.com/mailgun/mailgun-go"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

const (
	activationMail = `Hi {{.Name}},

Welcome to maji.moe. Please click the following link to activate your account:

http://maji.moe/activation/{{.ActivationToken}}

Thanks,
maji.moe`
	mailSender      = "maji.moe <noreply@maji.moe>"
	activationTitle = "Activate your account"
	mailRecipient   = "%s <%s>"
)

var activationTmpl = template.Must(template.New("activation mail").Parse(activationMail))

type EmailResendForm struct {
	Email int64 `form:"email"`
}

func (form *EmailResendForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Email, "email").Required("").Email("")

	return errors
}

func sendActivationMail(user *models.User, mg mailgun.Mailgun) {
	if user.Activated {
		return
	}

	var buf bytes.Buffer
	activationTmpl.Execute(&buf, user)

	recipient := fmt.Sprintf(mailRecipient, user.Name, user.Email)
	msg := mailgun.NewMessage(mailSender, activationTitle, buf.String(), recipient)

	if _, _, err := mg.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func EmailResend(form EmailResendForm, db *gorp.DbMap, r render.Render, mg mailgun.Mailgun) {
	var user models.User

	if err := db.SelectOne(&user, "SELECT * FROM users WHERE email=?", form.Email); err != nil {
		errors := NewErr([]string{"email"}, "213", "User does not exist")
		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	if user.Activated {
		errors := NewErr([]string{"email"}, "215", "User has been activated")
		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	go sendActivationMail(&user, mg)
}
