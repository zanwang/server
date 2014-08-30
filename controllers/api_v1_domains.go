package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/util"
)

func (a *APIv1) DomainList(c *gin.Context) {
	var domains []models.Domain
	user := c.MustGet("user").(*models.User)

	if _, err := models.DB.Select(&domains, "SELECT * FROM domains WHERE user_id=?", user.ID); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, domains)
}

type domainForm struct {
	Name *string
}

func (f *domainForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Name: "name",
	}
}

func (a *APIv1) DomainCreate(c *gin.Context) {
	token := c.MustGet("token").(*models.Token)
	user := c.MustGet("user").(*models.User)

	if token.UserID != user.ID {
		panic(errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.DomainForbidden,
			Message: "You are forbidden to create domains for this user",
		})
	}

	if !user.Activated {
		panic(errors.API{
			Status:  http.StatusForbidden,
			Code:    errors.UserNotActivated,
			Message: "User has not been activated",
		})
	}

	var form domainForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Name == nil {
		panic(errors.New("name", errors.Required, "Name is required"))
	}

	domain := models.Domain{
		Name:   *form.Name,
		UserID: token.UserID,
	}

	if err := models.DB.Insert(&domain); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusCreated, domain)
}

func (a *APIv1) DomainShow(c *gin.Context) {
	domain := c.MustGet("domain").(*models.Domain)

	util.Render.JSON(c.Writer, http.StatusOK, domain)
}

func (a *APIv1) DomainUpdate(c *gin.Context) {
	var form domainForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	domain := c.MustGet("domain").(*models.Domain)

	if form.Name != nil {
		domain.Name = *form.Name
	}

	if _, err := models.DB.Update(domain); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, domain)
}

func (a *APIv1) DomainDestroy(c *gin.Context) {
	domain := c.MustGet("domain").(*models.Domain)

	if _, err := models.DB.Delete(domain); err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (a *APIv1) DomainRenew(c *gin.Context) {
	domain := c.MustGet("domain").(*models.Domain)
	domain.Renew()

	if _, err := models.DB.Update(domain); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, domain)
}
