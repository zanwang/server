package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
)

type Middleware struct{}

func (m *Middleware) GetToken(c *gin.Context) {
	var token models.Token
	auth := strings.Split(c.Request.Header.Get("Authorization"), " ")

	if strings.ToLower(auth[0]) != "token" || auth[1] == "" {
		return
	}

	if err := models.DB.SelectOne(&token, "SELECT * FROM tokens WHERE key=?", auth[1]); err != nil {
		return
	}

	if token.ExpiredAt.Before(time.Now()) {
		if _, err := models.DB.Delete(&token); err != nil {
			panic(err)
		}

		panic(errors.API{
			Status:  http.StatusUnauthorized,
			Code:    errors.TokenExpired,
			Message: "Token is expired",
		})
	}

	c.Set("token", &token)
}

func (m *Middleware) TokenRequired(c *gin.Context) {
	m.GetToken(c)

	if _, err := c.Get("token"); err != nil {
		panic(errors.API{
			Status:  http.StatusUnauthorized,
			Code:    errors.Unauthorized,
			Message: "Token is required",
		})
	}
}

func (m *Middleware) GetUser(c *gin.Context) {
	var user models.User
	userID := c.Params.ByName("user_id")

	if err := models.DB.SelectOne(&user, "SELECT * FROM users WHERE id=?", userID); err != nil {
		panic(errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.NotFound,
			Message: "User does not exist",
		})
	}

	c.Set("user", &user)
}

func (m *Middleware) GetDomain(c *gin.Context) {
	var domain models.Domain
	domainID := c.Params.ByName("domain_id")

	if err := models.DB.SelectOne(&domain, "SELECT * FROM domains WHERE id=?", domainID); err != nil {
		panic(errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.NotFound,
			Message: "Domain does not exist",
		})
	}

	c.Set("domain", &domain)
}

func (m *Middleware) CheckOwnershipOfDomain(c *gin.Context) {
	m.GetDomain(c)

	token := c.MustGet("token").(*models.Token)
	domain := c.MustGet("domain").(*models.Domain)

	if domain.UserID != token.UserID {
		panic(errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.Forbidden,
			Message: "You are forbidden to access this domain",
		})
	}
}

func (m *Middleware) GetRecord(c *gin.Context) {
	var record models.Record
	var domain models.Domain
	recordID := c.Params.ByName("record_id")

	if err := models.DB.SelectOne(&record, "SELECT * FROM records WHERE id=?", recordID); err != nil {
		panic(errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.NotFound,
			Message: "Record does not exist",
		})
	}

	if err := models.DB.SelectOne(&domain, "SELECT * FROM domains WHERE id=?", record.DomainID); err != nil {
		panic(errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.NotFound,
			Message: "Domain does not exist",
		})
	}

	c.Set("record", &record)
	c.Set("domain", &domain)
}

func (m *Middleware) CheckOwnershipOfRecord(c *gin.Context) {
	m.GetRecord(c)

	token := c.MustGet("token").(*models.Token)
	domain := c.MustGet("domain").(*models.Domain)

	if token.UserID != domain.UserID {
		panic(errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.Forbidden,
			Message: "You are forbidden to access this record",
		})
	}
}
