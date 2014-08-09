package middlewares

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/models"
)

// GetCurrentUser gets the current user and map it to the context
func GetCurrentUser(c martini.Context, s sessions.Session, dbMap *gorp.DbMap) {
	var user models.User

	if err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE id=?", s.Get("userId")); err != nil {
		user.LoggedIn = false
		s.Delete("userId")
	} else {
		user.LoggedIn = true
	}

	c.Map(&user)
}

// NeedLogin checks whether the user has logged in
func NeedLogin(user *models.User, r render.Render) {
	if !user.LoggedIn {
		r.Redirect("/login")
	}
}
