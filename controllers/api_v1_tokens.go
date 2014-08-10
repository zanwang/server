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

func TokenCreate(form TokenCreateForm, r render.Render, dbMap *gorp.DbMap, token *models.Token) {
	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT id, password FROM users WHERE name=? OR email=?", form.Login, form.Login); err != nil {
		r.Status(http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		r.Status(http.StatusUnauthorized)
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

	token = &models.Token{
		UserID: user.ID,
		Key:    uniuri.New(),
	}

	if err := dbMap.Insert(token); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, token)
}

func TokenDestroy(dbMap *gorp.DbMap, res http.ResponseWriter, token *models.Token) {
	if count, err := dbMap.Delete(token); count > 0 {
		res.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		panic(err)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
