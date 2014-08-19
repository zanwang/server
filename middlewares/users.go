package middlewares

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/controllers"
	"github.com/tommy351/maji.moe/models"
)

func GetUser(c martini.Context, params martini.Params, db *gorp.DbMap, res http.ResponseWriter, token *models.Token) {
	var user models.User

	if err := db.SelectOne(&user, "SELECT * FROM users WHERE id=?", params["user_id"]); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	} else if user.ID == token.UserID {
		user.LoggedIn = true
	}

	c.Map(&user)
}

func NeedActivation(user *models.User, r render.Render) {
	if !user.Activated {
		errors := controllers.NewErr([]string{"common"}, "210", "User has not been activated")
		r.JSON(http.StatusForbidden, controllers.FormatErr(errors))
		return
	}
}
