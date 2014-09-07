package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
)

type Middleware struct{}

func (m *Middleware) GetToken(c *gin.Context) {
	auth := strings.Split(c.Request.Header.Get("Authorization"), " ")

	if strings.ToLower(auth[0]) != "token" || auth[1] == "" {
		return
	}

	var token models.Token

	if err := models.DB.Where("`key` = ?", auth[1]).Find(&token).Error; err != nil {
		return
	}

	if token.IsExpired() {
		if err := models.DB.Delete(&token).Error; err != nil {
			log.Println(err)
		}

		panic(&errors.API{
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
		panic(&errors.API{
			Status:  http.StatusUnauthorized,
			Code:    errors.TokenRequired,
			Message: "Token is required",
		})
	}
}

func (m *Middleware) GetUser(c *gin.Context) {
	var user models.User

	if err := models.DB.First(&user, c.Params.ByName("user_id")).Error; err != nil {
		panic(&errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.UserNotExist,
			Message: "User does not exist",
		})
	}

	c.Set("user", &user)
}

func (m *Middleware) CheckPermissionOfUser(c *gin.Context) {
	m.GetUser(c)

	var token *models.Token

	if data, err := c.Get("token"); err != nil {
		panic(&errors.API{
			Status:  http.StatusUnauthorized,
			Code:    errors.TokenRequired,
			Message: "Token is required",
		})
	} else {
		token = data.(*models.Token)
	}

	user := c.MustGet("user").(*models.User)

	if token.UserId != user.Id {
		panic(&errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.UserForbidden,
			Message: "You are forbidden to access this user",
		})
	}
}

func (m *Middleware) GetDomain(c *gin.Context) {
	var domain models.Domain

	if err := models.DB.First(&domain, c.Params.ByName("domain_id")).Error; err != nil {
		panic(&errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.DomainNotExist,
			Message: "Domain does not exist",
		})
	}

	c.Set("domain", &domain)
}

func (m *Middleware) CheckOwnershipOfDomain(c *gin.Context) {
	m.GetDomain(c)

	var token *models.Token

	if data, err := c.Get("token"); err != nil {
		panic(&errors.API{
			Status:  http.StatusUnauthorized,
			Code:    errors.TokenRequired,
			Message: "Token is required",
		})
	} else {
		token = data.(*models.Token)
	}

	domain := c.MustGet("domain").(*models.Domain)

	if domain.UserId != token.UserId {
		panic(&errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.DomainForbidden,
			Message: "You are forbidden to access this domain",
		})
	}
}

func (m *Middleware) GetRecord(c *gin.Context) {
	var record models.Record

	if err := models.DB.First(&record, c.Params.ByName("record_id")).Error; err != nil {
		panic(&errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.RecordNotExist,
			Message: "Record does not exist",
		})
	}

	var domain models.Domain
	if err := models.DB.First(&domain, record.DomainId).Error; err != nil {
		panic(&errors.API{
			Status:  http.StatusNotFound,
			Code:    errors.DomainNotExist,
			Message: "Domain does not exist",
		})
	}

	c.Set("record", &record)
	c.Set("domain", &domain)
}

func (m *Middleware) CheckOwnershipOfRecord(c *gin.Context) {
	m.GetRecord(c)

	var token *models.Token

	if data, err := c.Get("token"); err != nil {
		panic(&errors.API{
			Status:  http.StatusUnauthorized,
			Code:    errors.TokenRequired,
			Message: "Token is required",
		})
	} else {
		token = data.(*models.Token)
	}

	domain := c.MustGet("domain").(*models.Domain)

	if token.UserId != domain.UserId {
		panic(&errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.RecordForbidden,
			Message: "You are forbidden to access this record",
		})
	}
}
