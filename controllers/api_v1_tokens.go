package controllers

import (
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/go-martini/martini"
	"github.com/huandu/facebook"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

type TokenCreateForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (form *TokenCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Email, "email").Required("").Email("")
	v.Validate(&form.Password, "password").Required("").Length(6, 50, "")

	return errors
}

func TokenCreate(form TokenCreateForm, r render.Render, db *gorp.DbMap, c martini.Context) {
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
	}

	if err := db.Insert(&token); err != nil {
		panic(err)
	}

	c.Map(&token)
}

func TokenUpdate(db *gorp.DbMap, token *models.Token, res http.ResponseWriter, c martini.Context) {
	if count, err := db.Update(token); err != nil {
		panic(err)
	} else if count == 0 {
		res.WriteHeader(http.StatusNotFound)
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

type TokenFacebookForm struct {
	UserID      string `form:"user_id" json:"user_id"`
	AccessToken string `form:"access_token" json:"access_token"`
}

func (form *TokenFacebookForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.UserID, "user_id").Required("")
	v.Validate(&form.AccessToken, "access_token").Required("")

	return errors
}

func TokenFacebook(form TokenFacebookForm, db *gorp.DbMap, fb *facebook.App, r render.Render, c martini.Context) {
	session := fb.Session(form.AccessToken)
	res, err := session.Get("/"+form.UserID, nil)

	if err != nil {
		errors := NewErr([]string{"access_token"}, "216", "Facebook login failed")
		r.JSON(http.StatusBadRequest, errors)
		return
	}

	var user models.User
	id := res["id"].(string)
	name := res["name"].(string)
	email := res["email"].(string)

	if err := db.SelectOne(&user, "SELECT id FROM users WHERE facebook_id=?", id); err == nil {
		token := models.Token{
			UserID: user.ID,
		}

		if err := db.Insert(&token); err != nil {
			panic(err)
		}

		c.Map(&token)
		return
	}

	user = models.User{
		Name:       name,
		Email:      email,
		Avatar:     "//graph.facebook.com/" + form.UserID + "/picture",
		Activated:  true,
		FacebookID: id,
	}

	if err := db.Insert(&user); err != nil {
		if err.Error() == models.EmailTakenError {
			errors := NewErr([]string{"email"}, "211", "Email has been taken")
			r.JSON(http.StatusBadRequest, FormatErr(errors))
			return
		}

		panic(err)
	}

	token := models.Token{
		UserID: user.ID,
		Key:    uniuri.NewLen(32),
	}

	if err := db.Insert(&token); err != nil {
		panic(err)
	}

	c.Map(&token)
}

func TokenTwitter() {
	//
}

func TokenGoogle() {
	//
}

func TokenGitHub() {
	//
}
