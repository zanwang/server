package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
	"github.com/mholt/binding"
)

type userForm struct {
	Name        *string `json:"name"`
	OldPassword *string `json:"old_password"`
	Password    *string `json:"password"`
	Email       *string `json:"email"`
}

func (f *userForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Name:        "name",
		&f.OldPassword: "old_password",
		&f.Password:    "password",
		&f.Email:       "email",
	}
}

func handleUserDBError(err error) {
	switch e := err.(type) {
	case *mysql.MySQLError:
		switch e.Number {
		case errors.MySQLDuplicateEntry:
			panic(&errors.API{
				Field:   "email",
				Code:    errors.EmailUsed,
				Message: "Email has been taken",
			})
		default:
			panic(e)
		}
	default:
		panic(e)
	}
}

func (a *APIv1) UserCreate(c *gin.Context) {
	var form userForm

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

	user.SetActivated(false)
	user.Gravatar()

	if err := models.DB.Create(&user).Error; err != nil {
		handleUserDBError(err)
	}

	util.Render.JSON(c.Writer, http.StatusCreated, user)
	go user.SendActivationMail()
}

func (a *APIv1) UserShow(c *gin.Context) {
	var isOwner bool
	user := c.MustGet("user").(*models.User)

	if data, err := c.Get("token"); err == nil {
		token := data.(*models.Token)
		isOwner = user.Id == token.UserId
	}

	if isOwner {
		util.Render.JSON(c.Writer, http.StatusOK, user)
	} else {
		util.Render.JSON(c.Writer, http.StatusOK, user.PublicProfile())
	}
}

func (a *APIv1) UserUpdate(c *gin.Context) {
	var form userForm
	user := c.MustGet("user").(*models.User)

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Name != nil {
		user.Name = *form.Name
	}

	if form.Password != nil {
		if len(user.Password) > 0 {
			if form.OldPassword == nil {
				panic(errors.New("old_password", errors.Required, "Current password is required"))
			}

			if err := user.Authenticate(*form.OldPassword); err != nil {
				if e, ok := err.(*errors.API); ok {
					e.Field = "old_password"

					if e.Status == http.StatusUnauthorized {
						e.Status = http.StatusForbidden
					}

					panic(e)
				} else {
					panic(err)
				}
			}
		}

		if err := user.GeneratePassword(*form.Password); err != nil {
			panic(err)
		}
	}

	if form.Email != nil && user.Email != *form.Email {
		user.Email = *form.Email
		user.Gravatar()
		user.SetActivated(false)
	}

	if err := models.DB.Save(user).Error; err != nil {
		handleUserDBError(err)
	}

	util.Render.JSON(c.Writer, http.StatusOK, user)
	go user.SendActivationMail()
}

func (a *APIv1) UserDestroy(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	if err := models.DB.Delete(user).Error; err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
