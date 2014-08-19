package middlewares

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/tommy351/maji.moe/models"
)

func GetRecord(params martini.Params, c martini.Context, db *gorp.DbMap, res http.ResponseWriter) {
	var record models.Record

	if err := db.SelectOne(&record, "SELECT * FROM records WHERE id=?", params["record_id"]); err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	c.Map(&record)
}

func CheckOwnershipOfRecord(strict bool) martini.Handler {
	return func(record *models.Record, db *gorp.DbMap, token *models.Token, res http.ResponseWriter, c martini.Context) {
		var domain models.Domain

		if err := db.SelectOne(&domain, "SELECT * FROM domains WHERE id=?", record.DomainID); err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if domain.UserID != token.UserID && (strict || !domain.Public) {
			res.WriteHeader(http.StatusForbidden)
			return
		}

		c.Map(&domain)
	}
}
