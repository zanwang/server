package controllers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
	"github.com/mholt/binding"
)

func (a *APIv1) PasswordReset(c *gin.Context) {
	var form emailForm

	if err := binding.Bind(c.Request, &form); err != nil {
		bindingError(err)
	}

	if form.Email == nil {
		panic(errors.New("email", errors.Required, "Email is required"))
	}

	if !govalidator.IsEmail(*form.Email) {
		panic(errors.New("email", errors.Email, "Email is invalid"))
	}

	var user models.User

	if err := models.DB.Where("email = ?", *form.Email).Find(&user).Error; err != nil {
		panic(&errors.API{
			Status:  http.StatusNotFound,
			Field:   "email",
			Code:    errors.UserNotExist,
			Message: "User does not exist",
		})
	}

	user.PasswordResetToken = uniuri.NewLen(32)

	if err := models.DB.Save(&user).Error; err != nil {
		panic(err)
	}

	util.Render.JSON(c.Writer, http.StatusAccepted, map[string]interface{}{
		"email": user.Email,
	})

	go user.SendPasswordResetMail()
}
