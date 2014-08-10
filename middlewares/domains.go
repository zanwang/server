package middlewares

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/tommy351/maji.moe/models"
)

func GetDomain(params martini.Params, c martini.Context, db *gorp.DbMap, res http.ResponseWriter) {
	var domain models.Domain

	if err := db.SelectOne(&domain, "SELECT * FROM domains WHERE id=?", params["domain_id"]); err != nil {
		res.WriteHeader(http.StatusNotFound)
	}

	c.Map(&domain)
}

func CheckOwnershipOfDomain(strict bool) martini.Handler {
	return func(token *models.Token, res http.ResponseWriter, domain *models.Domain) {
		if domain.UserID != token.UserID && (strict || !domain.Public) {
			res.WriteHeader(http.StatusForbidden)
		}
	}
}
