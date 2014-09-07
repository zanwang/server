package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/models"
	"github.com/majimoe/server/util"
)

func Home(c *gin.Context) {
	util.Render.HTML(c.Writer, http.StatusOK, "index", nil)
}

func App(c *gin.Context) {
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
		util.Render.HTML(c.Writer, http.StatusBadRequest, "error/activation", nil)
		return
	}

	user.Activated = true

	if err := models.DB.Save(&user).Error; err != nil {
		panic(err)
	}

	http.Redirect(c.Writer, c.Request, "/app", http.StatusFound)
}
