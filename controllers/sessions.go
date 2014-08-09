package controllers

import (
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/middlewares"
	"github.com/tommy351/maji.moe/models"
)

// SessionNew handles GET /login
func SessionNew(r render.Render, currentUser *models.User, csrf *middlewares.CSRFToken, s sessions.Session) {
	if currentUser.LoggedIn {
		r.Redirect("/app")
		return
	}

	var errors map[string]interface{}

	if flashes := s.Flashes(); len(flashes) > 0 {
		errors = formatErr(flashes[0])
	}

	r.HTML(http.StatusOK, "sessions/new", map[string]interface{}{
		"Errors": errors,
		"Token":  csrf.GetToken(),
	})
}

// SessionCreateForm binds data for SessionCreate
type SessionCreateForm struct {
	Login    string `form:"login" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// SessionCreate handles POST /sessions
func SessionCreate(form SessionCreateForm, r render.Render, s sessions.Session, errors binding.Errors, dbMap *gorp.DbMap) {
	if errors != nil {
		s.AddFlash(errors)
		r.Redirect("/login")
		return
	}

	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE username=? OR email=?", form.Login, form.Login); err != nil {
		s.AddFlash(map[string]interface{}{
			"login": "User does not exist",
		})
		r.Redirect("/login")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		s.AddFlash(map[string]interface{}{
			"password": "Wrong password",
		})
		r.Redirect("/login")
		return
	}

	s.Set("userId", user.ID)
	r.Redirect("/app")
}

// SessionDestroy handles GET /logout & DELETE /sessions
func SessionDestroy(s sessions.Session, r render.Render) {
	s.Delete("userId")
	r.Redirect("/login")
}
