package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/errors"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/util"
)

type userForm struct {
	Name        *string
	OldPassword *string
	Password    *string
	Email       *string
}

func (f *userForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Name:        "name",
		&f.OldPassword: "old_password",
		&f.Password:    "password",
		&f.Email:       "email",
	}
}

func (a *APIv1) UserCreate(c *gin.Context) {
	var form userForm
	conf := config.Config

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Name == nil {
		panic(errors.New("name", errors.Required, "Name is required"))
	}

	if form.Password == nil {
		panic(errors.New("password", errors.Required, "Password is required"))
	}

	if form.Email == nil {
		panic(errors.New("email", errors.Required, "Email is required"))
	}

	user := models.User{
		Name:  *form.Name,
		Email: *form.Email,
	}

	if err := user.GeneratePassword(*form.Password); err != nil {
		panic(err)
	}

	user.SetActivated(!conf.EmailActivation)
	user.Gravatar()

	if err := models.DB.Insert(&user); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusCreated, user)
	go user.SendActivationMail()
}

func (a *APIv1) UserShow(c *gin.Context) {
	var isOwner bool
	user := c.MustGet("user").(*models.User)

	if token, ok := c.MustGet("token").(*models.Token); ok {
		isOwner = user.ID == token.UserID
	}

	if isOwner {
		util.Render.JSON(c.Writer, http.StatusOK, user)
	} else {
		util.Render.JSON(c.Writer, http.StatusOK, map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"avatar":     user.Avatar,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}
}

func (a *APIv1) UserUpdate(c *gin.Context) {
	var form userForm
	conf := config.Config

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	user := c.MustGet("user").(*models.User)

	if form.Name != nil {
		user.Name = *form.Name
	}

	if form.Password != nil {
		if form.OldPassword == nil {
			panic(errors.New("old_password", errors.Required, "Current password is required"))
		}

		if len(user.Password) > 0 {
			if err := user.Authenticate(*form.OldPassword); err != nil {
				panic(err)
			}
		}

		if err := user.GeneratePassword(*form.Password); err != nil {
			panic(err)
		}
	}

	if form.Email != nil && user.Email != *form.Email {
		user.Email = *form.Email
		user.Gravatar()

		if conf.EmailActivation {
			user.SetActivated(false)
		}
	}

	if _, err := models.DB.Update(user); err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, user)
	go user.SendActivationMail()
}

func (a *APIv1) UserDestroy(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	if _, err := models.DB.Delete(user); err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
