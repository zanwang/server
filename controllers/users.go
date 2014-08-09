package controllers

import (
	"net/http"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/middlewares"
	"github.com/tommy351/maji.moe/models"
)

// UserNew handles GET /signup
func UserNew(r render.Render, currentUser *models.User, csrf *middlewares.CSRFToken, s sessions.Session) {
	if currentUser.LoggedIn {
		r.Redirect("/users/" + currentUser.Username)
		return
	}

	var errors map[string]interface{}

	if flashes := s.Flashes(); len(flashes) > 0 {
		errors = formatErr(flashes[0])
	}

	r.HTML(http.StatusOK, "users/new", map[string]interface{}{
		"Errors": errors,
		"Token":  csrf.GetToken(),
	})
}

// UserCreateForm binds data for UserCreate
type UserCreateForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Confirm  string `form:"confirm" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

// Validate validates UserCreateForm
func (form *UserCreateForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := Validation{Errors: &errors}

	v.Validate(&form.Username, "username").Length(3, 20, "")
	v.Validate(&form.Password, "password").Length(6, 50, "")
	v.Validate(&form.Confirm, "confirm").Equal(&form.Password, "")
	v.Validate(&form.Email, "email").Email("")

	return errors
}

// UserCreate handles POST /users
func UserCreate(form UserCreateForm, r render.Render, errors binding.Errors, s sessions.Session, dbMap *gorp.DbMap) {
	if errors != nil {
		s.AddFlash(errors)
		r.Redirect("/signup")
		return
	}

	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE username=? OR email=?", form.Username, form.Email); err == nil {
		errors = binding.Errors{}

		if form.Username == user.Username {
			errors.Add([]string{"username"}, "ContentError", "Username has been used")
		}

		if form.Email == user.Email {
			errors.Add([]string{"email"}, "ContentError", "Email has been used")
		}

		s.AddFlash(errors)
		r.Redirect("/signup")
		return
	}

	password, err := generatePassword(form.Password)
	if err != nil {
		s.AddFlash(map[string]interface{}{
			"common": "Internal server error",
		})
		r.Redirect("/signup")
		return
	}

	now := time.Now().UnixNano()

	user = models.User{
		Username:        form.Username,
		Password:        string(password),
		Email:           form.Email,
		CreatedAt:       now,
		UpdatedAt:       now,
		Activated:       false,
		ActivationToken: uniuri.New(),
	}

	if err := dbMap.Insert(&user); err != nil {
		s.AddFlash(map[string]interface{}{
			"common": "Internal server error",
		})
		r.Redirect("/signup")
		return
	}

	s.Set("userId", user.ID)
	r.Redirect("/app")
}

// UserEdit handles GET /settings/profile
func UserEdit() {
	//
}

// UserUpdate handles PUT /users/:user_id
func UserUpdate() {
	//
}

// UserDestroy handles DELETE /users/:user_id
func UserDestroy() {
	//
}
