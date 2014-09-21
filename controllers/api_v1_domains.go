package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
	"github.com/mholt/binding"
)

func (a *APIv1) DomainList(c *gin.Context) {
	domains := []models.Domain{}
	user := c.MustGet("user").(*models.User)

	if err := models.DB.Where("user_id = ?", user.Id).Find(&domains).Error; err != nil {
		if err != gorm.RecordNotFound {
			panic(err)
		}
	}

	util.Render.JSON(c.Writer, http.StatusOK, domains)
}

type domainForm struct {
	Name *string `json:"name"`
}

func (f *domainForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Name: "name",
	}
}

func handleDomainDBError(err error) {
	switch e := err.(type) {
	case *mysql.MySQLError:
		switch e.Number {
		case errors.MySQLDuplicateEntry:
			panic(&errors.API{
				Field:   "name",
				Code:    errors.DomainUsed,
				Message: "Domain name has been taken",
			})
		default:
			panic(e)
		}
	default:
		panic(e)
	}
}

func (a *APIv1) DomainCreate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	if !user.Activated {
		panic(&errors.API{
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
		UserId: user.Id,
	}

	if err := models.DB.Create(&domain).Error; err != nil {
		handleDomainDBError(err)
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

	if err := models.DB.Save(domain).Error; err != nil {
		handleDomainDBError(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, domain)
}

func (a *APIv1) DomainDestroy(c *gin.Context) {
	domain := c.MustGet("domain").(*models.Domain)

	if err := models.DB.Delete(domain).Error; err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (a *APIv1) DomainRenew(c *gin.Context) {
	domain := c.MustGet("domain").(*models.Domain)

	if err := domain.Renew(); err != nil {
		panic(err)
	}

	if err := models.DB.Save(domain).Error; err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, domain)
}
