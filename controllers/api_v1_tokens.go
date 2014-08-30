package controllers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/util"
)

type tokenForm struct {
	Email    *string
	Password *string
}

func (f *tokenForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Email:    "email",
		&f.Password: "password",
	}
}

func responseToken(c *gin.Context, token *models.Token) {
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Expires", "0")
	util.Render.JSON(c.Writer, http.StatusCreated, token)
}

func (a *APIv1) TokenCreate(c *gin.Context) {
	var user models.User
	var form tokenForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Email == nil {
		panic(errors.New("email", errors.Required, "Email is required"))
	}

	if form.Password == nil {
		panic(errors.New("password", errors.Required, "Password is required"))
	}

	if !govalidator.IsEmail(*form.Email) {
		panic(errors.New("email", errors.Email, "Email is invalid"))
	}

	if err := models.DB.SelectOne(&user, "SELECT id, password FROM users WHERE email=?", *form.Email); err != nil {
		panic(errors.API{
			Status:  http.StatusBadRequest,
			Field:   "email",
			Code:    errors.UserNotExist,
			Message: "User does not exist",
		})
	}

	if err := user.Authenticate(*form.Password); err != nil {
		panic(err)
	}

	token := models.Token{UserID: user.ID}

	if err := models.DB.Insert(&token); err != nil {
		panic(err)
	}

	responseToken(c, &token)
}

func (a *APIv1) TokenUpdate(c *gin.Context) {
	token := c.MustGet("token").(*models.Token)

	if _, err := models.DB.Update(token); err != nil {
		panic(err)
	}

	responseToken(c, token)
}

func (a *APIv1) TokenDestroy(c *gin.Context) {
	token := c.MustGet("token").(*models.Token)

	if _, err := models.DB.Delete(token); err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
