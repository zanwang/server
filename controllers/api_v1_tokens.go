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

func TokenCreate(form TokenCreateForm, r render.Render, db *gorp.DbMap) {
	var user models.User

	if err := db.SelectOne(&user, "SELECT id, password FROM users WHERE name=? OR email=?", form.Login, form.Login); err != nil {
		r.Status(http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		r.Status(http.StatusUnauthorized)
		return
	}

	var token models.Token

	if err := db.SelectOne(&token, "SELECT * FROM tokens WHERE user_id=?", user.ID); err != nil {
		// Create a token
		token = models.Token{
			UserID: user.ID,
			Key:    uniuri.NewLen(32),
		}

		if err := db.Insert(&token); err != nil {
			panic(err)
		}
	} else {
		// Update the existing token
		if _, err := db.Update(&token); err != nil {
			panic(err)
		}
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
