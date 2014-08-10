package controllers

import (
	"net/http"

	"github.com/martini-contrib/render"
)

// App is the entry point
func App(r render.Render) {
	r.HTML(http.StatusOK, "app", nil)
}

func Home(r render.Render) {
	r.HTML(http.StatusOK, "index", nil)
}
