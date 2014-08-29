package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tommy351/maji.moe/models"
	"github.com/tommy351/maji.moe/util"
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

func Activation(c *gin.Context) {
	var user models.User
	query := c.Request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")

	if err := models.DB.SelectOne(&user, "SELECT * FROM users WHERE id=?", id); err != nil {
		NotFound(c)
		return
	}

	if user.Activated {
		http.Redirect(c.Writer, c.Request, "/app", http.StatusFound)
		return
	}

	if user.ActivationToken != token {
		util.Render.HTML(c.Writer, http.StatusBadRequest, "error/activation", nil)
		return
	}

	user.Activated = true

	if _, err := models.DB.Update(&user); err != nil {
		panic(err)
	}

	http.Redirect(c.Writer, c.Request, "/app", http.StatusFound)
}
