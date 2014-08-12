package controllers

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/go-martini/martini"
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

	v.Validate(&form.Name, "name").Required("").Length(3, 20, "")
	v.Validate(&form.Password, "password").Required("").Length(6, 50, "")
	v.Validate(&form.Email, "email").Required("").Email("")

	return errors
}

func UserCreate(form UserCreateForm, r render.Render, dbMap *gorp.DbMap) {
	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT name,email FROM users WHERE name=? OR email=?", form.Name, form.Email); err == nil {
		errors := binding.Errors{}

		if form.Name == user.Name {
			errors.Add([]string{"name"}, "210", "User name has been taken")
		}

		if form.Email == user.Email {
			errors.Add([]string{"email"}, "211", "Email has been taken")
		}

		r.JSON(http.StatusBadRequest, FormatErr(errors))
		return
	}

	// TODO gravatar avatar
	user = models.User{
		Name:            form.Name,
		Password:        generatePassword(form.Password),
		Email:           form.Email,
		Activated:       false,
		ActivationToken: uniuri.New(),
	}

	user.Gravatar()

	if err := dbMap.Insert(&user); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, user)
}

func UserShow(r render.Render, user *models.User) {
	if user.LoggedIn {
		r.JSON(http.StatusOK, user)
	} else {
		r.JSON(http.StatusOK, map[string]interface{}{
			"id":           user.ID,
			"name":         user.Name,
			"display_name": user.DisplayName,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		})
	}
}

type UserUpdateForm struct {
	Name        string `form:"name"`
	DisplayName string `form:"display_name"`
	Password    string `form:"password"`
	Email       string `form:"email"`
}

func (form *UserUpdateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Name, "name").Length(3, 20, "")
	v.Validate(&form.Password, "password").Length(6, 50, "")
	v.Validate(&form.Email, "email").Email("")

	return errors
}

func UserUpdate(form UserUpdateForm, r render.Render, db *gorp.DbMap, user *models.User) {
	if form.Name != "" {
		user.Name = form.Name
	}

	user.DisplayName = form.DisplayName

	if form.Password != "" {
		user.Password = generatePassword(form.Password)
	}

	if form.Email != "" {
		user.Email = form.Email
	}

	if count, err := db.Update(user); count > 0 {
		r.JSON(http.StatusOK, user)
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
		r.Status(http.StatusNotFound)
		return
	}

	if user.Activated {
		r.Status(http.StatusBadRequest)
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
