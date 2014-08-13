package middlewares

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/tommy351/maji.moe/models"
)

func CheckToken(c martini.Context, db *gorp.DbMap, req *http.Request) {
	var token models.Token
	header := req.Header.Get("X-Auth-Token")

	if header != "" {
		if err := db.SelectOne(&token, "SELECT * FROM tokens WHERE key=?", header); err != nil {
			token.Authorized = true
		}
	}

	c.Map(&token)
	c.Next()

	// Update token
	if token.Authorized {
		if _, err := db.Update(&token); err != nil {
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
