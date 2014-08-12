package controllers

import (
	"net/http"

	"github.com/martini-contrib/binding"
)

type EmailResendForm struct {
	UserID int64 `form:"user_id"`
}

func (form *EmailResendForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.UserID, "user_id").Required("")

	return errors
}

func EmailResend(form EmailResendForm) {
	//
}
