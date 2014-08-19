package controllers

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/go-martini/martini"
	"github.com/mailgun/mailgun-go"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

type UserCreateForm struct {
	Name     string `form:"name" json:"name"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
}

func (form *UserCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Name, "name").Required("")
	v.Validate(&form.Password, "password").Required("").Length(6, 50, "")
	v.Validate(&form.Email, "email").Required("").Email("")

	return errors
}

func UserCreate(form UserCreateForm, r render.Render, dbMap *gorp.DbMap, mg mailgun.Mailgun) {
	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT email FROM users WHERE email=?", form.Email); err == nil {
		errors := NewErr([]string{"email"}, "211", "Email has been taken")
		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	user = models.User{
		Name:            form.Name,
		Password:        generatePassword(form.Password),
		Email:           form.Email,
		Activated:       false,
		ActivationToken: uniuri.NewLen(32),
		PasswordSet:     true,
	}

	user.Gravatar()

	if err := dbMap.Insert(&user); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, user)
	go sendActivationMail(&user, mg)
}

func UserShow(r render.Render, user *models.User) {
	if user.LoggedIn {
		r.JSON(http.StatusOK, user)
	} else {
		r.JSON(http.StatusOK, map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}
}

type UserUpdateForm struct {
	Name     string `form:"name"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

func (form *UserUpdateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Password, "password").Length(6, 50, "")
	v.Validate(&form.Email, "email").Email("")

	return errors
}

func UserUpdate(form UserUpdateForm, r render.Render, db *gorp.DbMap, user *models.User, mg mailgun.Mailgun) {
	if form.Name != "" {
		user.Name = form.Name
	}

	if form.Password != "" {
		user.Password = generatePassword(form.Password)
		user.PasswordSet = true
	}

	if form.Email != "" && user.Email != form.Email {
		user.Email = form.Email
		user.Activated = false
		user.ActivationToken = uniuri.NewLen(32)
	}

	if count, err := db.Update(user); count > 0 {
		r.JSON(http.StatusOK, user)
		go sendActivationMail(user, mg)
	} else {
		panic(err)
	}
}

func UserDestroy(res http.ResponseWriter, db *gorp.DbMap, user *models.User) {
	if count, err := db.Delete(user); count > 0 {
		res.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		panic(err)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

func UserActivate(params martini.Params, db *gorp.DbMap, r render.Render) {
	var user models.User

	if err := db.SelectOne(&user, "SELECT * FROM users WHERE activation_token=?", params["token"]); err != nil {
		r.HTML(http.StatusNotFound, "error/404", nil)
		return
	}

	if user.Activated {
		r.Redirect("/app")
		return
	}

	user.Activated = true

	if count, err := db.Update(&user); count > 0 {
		r.Redirect("/app")
	} else if err != nil {
		panic(err)
	} else {
		r.Status(http.StatusNotFound)
	}
}
