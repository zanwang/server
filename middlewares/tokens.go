package middlewares

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/tommy351/maji.moe/models"
)

func CheckToken(c martini.Context, dbMap *gorp.DbMap, req *http.Request) {
	var token models.Token
	header := req.Header.Get("X-Auth-Token")

	if err := dbMap.SelectOne(&token, "SELECT * FROM tokens WHERE key=?", header); err != nil {
		token.Authorized = true
	} else {
		token.Authorized = false
	}

	c.Map(&token)
	c.Next()

	// Update token
	if token.Authorized {
		if _, err := dbMap.Update(&token); err != nil {
			panic(err)
		}
	}
}

func NeedAuthorization(token *models.Token, res http.ResponseWriter) {
	if !token.Authorized {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

func CheckCurrentUser(user *models.User, res http.ResponseWriter) {
	if !user.LoggedIn {
		res.WriteHeader(http.StatusForbidden)
	}
}
