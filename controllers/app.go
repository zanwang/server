package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/errors"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
	"github.com/tommy351/gin-csrf"
)

func Home(c *gin.Context) {
	util.Render.HTML(c.Writer, http.StatusOK, "index", nil)
}

func App(c *gin.Context) {
	// Prevent app from being embedding
	c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")

	// CSRF token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "csrftoken",
		Value: csrf.GetToken(c),
	})

	util.Render.HTML(c.Writer, http.StatusOK, "app", nil)
}

func NotFound(c *gin.Context) {
	util.Render.HTML(c.Writer, http.StatusNotFound, "error/404", nil)
}

func UserActivation(c *gin.Context) {
	var user models.User

	if err := models.DB.First(&user, c.Params.ByName("user_id")).Error; err != nil {
		NotFound(c)
		return
	}

	if user.Activated {
		http.Redirect(c.Writer, c.Request, "/app", http.StatusFound)
		return
	}

	if user.ActivationToken != c.Params.ByName("token") {
		NotFound(c)
		return
	}

	user.Activated = true

	if err := models.DB.Save(&user).Error; err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/app")
}

func PasswordReset(c *gin.Context) {
	var user models.User

	if err := models.DB.First(&user, c.Params.ByName("user_id")).Error; err != nil {
		NotFound(c)
		return
	}

	if user.PasswordResetToken == "" || user.PasswordResetToken != c.Params.ByName("token") {
		NotFound(c)
		return
	}

	util.Render.HTML(c.Writer, http.StatusOK, "password_reset", map[string]interface{}{
		"user":  user,
		"csrf":  csrf.GetToken(c),
		"error": nil,
	})
}

func PasswordResetSubmit(c *gin.Context) {
	var user models.User
	var err error
	password := c.Request.FormValue("password")
	confirm := c.Request.FormValue("confirm")

	if err := models.DB.First(&user, c.Params.ByName("user_id")).Error; err != nil {
		NotFound(c)
		return
	}

	if user.PasswordResetToken == "" || user.PasswordResetToken != c.Params.ByName("token") {
		NotFound(c)
		return
	}

	if password == "" {
		err = errors.New("password", errors.Required, "Password is required")
	} else if confirm == "" {
		err = errors.New("confirm", errors.Required, "Password confirmation is required")
	} else if len(password) < 6 {
		err = errors.New("password", errors.MinLength, "Password must be at least 6 characters long")
	} else if len(password) > 50 {
		err = errors.New("password", errors.MaxLength, "Password should be no longer than 50 characters")
	} else if password != confirm {
		err = errors.New("confirm", errors.PasswordConfirm, "Password confirmation doesn't match")
	}

	if err != nil {
		util.Render.HTML(c.Writer, http.StatusBadRequest, "password_reset", map[string]interface{}{
			"user":  user,
			"csrf":  csrf.GetToken(c),
			"error": err,
		})

		return
	}

	if err := user.GeneratePassword(password); err != nil {
		panic(err)
	}

	user.PasswordResetToken = ""

	if err := models.DB.Save(&user).Error; err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/login")
}
