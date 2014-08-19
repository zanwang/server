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
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (form *TokenCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Email, "email").Required("").Email("")
	v.Validate(&form.Password, "password").Required("").Length(6, 50, "")

	return errors
}

func TokenCreate(form TokenCreateForm, r render.Render, db *gorp.DbMap) {
	var user models.User

	if err := db.SelectOne(&user, "SELECT id, password FROM users WHERE email=?", form.Email); err != nil {
		errors := NewErr([]string{"email"}, "213", "User does not exist")
		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		errors := NewErr([]string{"password"}, "214", "Password is wrong")
		r.JSON(http.StatusUnauthorized, FormatErr(errors))
		return
	}

	token := models.Token{
		UserID: user.ID,
		Key:    uniuri.NewLen(32),
	}

	if err := db.Insert(&token); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, token)
}

func TokenUpdate(db *gorp.DbMap, r render.Render, token *models.Token) {
	if count, err := db.Update(token); count > 0 {
		r.JSON(http.StatusOK, token)
	} else if err != nil {
		panic(err)
	} else {
		r.Status(http.StatusNotFound)
	}
}

func TokenDestroy(db *gorp.DbMap, res http.ResponseWriter, token *models.Token) {
	if count, err := db.Delete(token); count > 0 {
		res.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		panic(err)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
