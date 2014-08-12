package middlewares

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/tommy351/maji.moe/models"
)

func GetUser(c martini.Context, params martini.Params, db *gorp.DbMap, res http.ResponseWriter, token *models.Token) {
	var user models.User

	if err := db.SelectOne("SELECT * FROM users WHERE id=?", params["user_id"]); err != nil {
		res.WriteHeader(http.StatusNotFound)
	} else if user.ID == token.UserID {
		user.LoggedIn = true
	}

	c.Map(&user)
}

func NeedActivation(user *models.User, res http.ResponseWriter) {
	if !user.Activated {
		res.WriteHeader(http.StatusForbidden)
	}
}
