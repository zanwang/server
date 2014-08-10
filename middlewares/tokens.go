package middlewares

import (
	"net/http"
	"time"

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
		token.UpdatedAt = time.Now().UnixNano()
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

/*
func CheckCurrentUser(params martini.Params, token *models.Token, res http.ResponseWriter) {
	if userID, err := strconv.ParseInt(params["user_id"], 10, 64); err != nil || userID != token.UserID {
		res.WriteHeader(http.StatusForbidden)
	}
}
*/

func CheckCurrentUser(user *models.User, res http.ResponseWriter) {
	if !user.LoggedIn {
		res.WriteHeader(http.StatusForbidden)
	}
}
