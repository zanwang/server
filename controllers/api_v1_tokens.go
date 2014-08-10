package controllers

import (
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

type TokenCreateForm struct {
	Login    string `form:"login"`
	Password string `form:"password"`
}

func (form *TokenCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Login, "login").Required("")
	v.Validate(&form.Password, "password").Required("").Length(6, 50, "")

	return errors
}

func TokenCreate(form TokenCreateForm, r render.Render, errors binding.Errors, dbMap *gorp.DbMap, token *models.Token) {
	if errors != nil {
		r.JSON(http.StatusBadRequest, formatErr(errors))
		return
	}

	// Update existing token
	if token != nil {
		if _, err := dbMap.Update(&token); err != nil {
			panic(err)
		} else {
			r.JSON(http.StatusCreated, token)
		}

		return
	}

	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT id, password FROM users WHERE name=? OR email=?", form.Login, form.Login); err != nil {
		errors := newErr([]string{"common"}, "404", "User does not exist")
		r.JSON(http.StatusBadRequest, formatErr(errors))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		errors := newErr([]string{"common"}, "401", "Unauthorized")
		r.JSON(http.StatusUnauthorized, formatErr(errors))
		return
	}

	token = &models.Token{
		UserID: user.ID,
		Key:    uniuri.New(),
	}

	if err := dbMap.Insert(token); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, token)
}

func TokenDestroy(dbMap *gorp.DbMap, r render.Render, token *models.Token) {
	if count, err := dbMap.Delete(token); count > 0 {
		r.Status(http.StatusNoContent)
	} else {
		panic(err)
	}
}
