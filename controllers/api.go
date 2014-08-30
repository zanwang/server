package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/util"
)

type APIv1 struct{}

func (api *APIv1) Entry(c *gin.Context) {
	util.Render.JSON(c.Writer, http.StatusOK, map[string]interface{}{
		"tokens":  "/api/v1/tokens",
		"users":   "/api/v1/users",
		"domains": "/api/v1/domains",
		"records": "/api/v1/records",
		"emails":  "/api/v1/emails",
	})
}

func (api *APIv1) Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case errors.API:
				if e.Status == 0 {
					e.Status = http.StatusBadRequest
				} else if e.Status == http.StatusInternalServerError {
					showStack(err)
				}

				util.Render.JSON(c.Writer, e.Status, e)
			default:
				showStack(err)
				util.Render.JSON(c.Writer, http.StatusInternalServerError, errors.API{
					Code:    errors.ServerError,
					Message: "Server error",
				})
			}

			c.Abort(0)
		}
	}()

	c.Next()
}

func (api *APIv1) UpdateToken(c *gin.Context) {
	defer func() {
		if data, err := c.Get("token"); err == nil {
			token := data.(*models.Token)
			models.DB.Update(token)
		}
	}()

	c.Next()
}

func bindingError(err binding.Errors) {
	if len(err) == 0 {
		return
	}

	var code int

	switch err[0].Classification {
	case binding.RequiredError:
		code = errors.Required
	case binding.ContentTypeError:
		code = errors.ContentType
	case binding.DeserializationError:
		code = errors.Deserialization
	case binding.TypeError:
		code = errors.Type
	default:
		code = errors.Unknown
	}

	panic(errors.API{Code: code, Message: err[0].Message})
}
