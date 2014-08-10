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

func UserCreate(form UserCreateForm, r render.Render, errors binding.Errors, dbMap *gorp.DbMap) {
	if errors != nil {
		r.JSON(http.StatusBadRequest, formatErr(errors))
		return
	}

	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE name=? OR email=?", form.Name, form.Email); err == nil {
		errors = binding.Errors{}

		if form.Name == user.Name {
			errors.Add([]string{"name"}, "210", "User name has been taken")
		}

		if form.Email == user.Email {
			errors.Add([]string{"email"}, "211", "Email has been taken")
		}

		r.JSON(http.StatusBadRequest, formatErr(errors))
		return
	}

	user = models.User{
		Name:            form.Name,
		Password:        generatePassword(form.Password),
		Email:           form.Email,
		Activated:       false,
		ActivationToken: uniuri.New(),
	}

	if err := dbMap.Insert(&user); err != nil {
		panic(err)
	}

	r.JSON(http.StatusCreated, user)
}

func UserShow(params martini.Params, r render.Render, dbMap *gorp.DbMap, token *models.Token) {
	var user models.User
	userID := toint64(params["user_id"])

	if err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE id=?", userID); err != nil {
		errors := newErr([]string{"common"}, "404", "User does not exist")
		r.JSON(http.StatusNotFound, formatErr(errors))
		return
	}

	if token.UserID == userID {
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

func UserUpdate(form UserUpdateForm, errors binding.Errors, r render.Render, dbMap *gorp.DbMap, params martini.Params, token *models.Token) {
	if errors != nil {
		r.JSON(http.StatusBadRequest, formatErr(errors))
		return
	}

	var user models.User
	userID := toint64(params["user_id"])

	if err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE id=?", userID); err != nil {
		errors := newErr([]string{"common"}, "404", "User does not exist")
		r.JSON(http.StatusBadRequest, formatErr(errors))
		return
	}

	if userID != token.UserID {
		errors := newErr([]string{"common"}, "403", "Forbidden")
		r.JSON(http.StatusForbidden, formatErr(errors))
		return
	}

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

	if count, err := dbMap.Update(&user); count > 0 {
		r.JSON(http.StatusOK, user)
	} else {
		panic(err)
	}
}

func UserDestroy(r render.Render, dbMap *gorp.DbMap, params martini.Params, token *models.Token) {
	userID := toint64(params["user_id"])

	if userID != token.UserID {
		errors := newErr([]string{"common"}, "403", "Forbidden")
		r.JSON(http.StatusForbidden, formatErr(errors))
		return
	}

	user := models.User{ID: userID}

	if count, err := dbMap.Delete(&user); count > 0 {
		r.Status(http.StatusNoContent)
	} else {
		panic(err)
	}
}
