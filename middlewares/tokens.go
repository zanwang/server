package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/models"
)

func CheckToken(c martini.Context, db *gorp.DbMap, req *http.Request) {
	var token models.Token
	auth := strings.Split(req.Header.Get("Authorization"), " ")

	if strings.ToLower(auth[0]) == "token" && auth[1] != "" {
		if err := db.SelectOne(&token, "SELECT * FROM tokens WHERE key=?", auth[1]); err == nil {
			// Delete token if expired
			if token.ExpiredAt.Before(time.Now()) {
				if _, err := db.Delete(&token); err != nil {
					panic(err)
				}
			} else {
				token.Authorized = true
			}
		}
	}

	c.Map(&token)

	// Update token
	if token.Authorized {
		defer func() {
			if _, err := db.Update(&token); err != nil {
				panic(err)
			}
		}()
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

func ResponseToken(token *models.Token, res http.ResponseWriter, r render.Render) {
	res.Header().Set("Pragma", "no-cache")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	res.Header().Set("Expires", "0")
	r.JSON(http.StatusOK, token)
}
