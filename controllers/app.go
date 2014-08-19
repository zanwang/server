package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/martini-contrib/render"
)

// App is the entry point
func App(r render.Render) {
	r.HTML(http.StatusOK, "app", nil)
}

func Home(r render.Render) {
	log.Print("index")
	r.HTML(http.StatusOK, "index", nil)
}

func NotFound(r render.Render, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api") {
		errors := NewErr([]string{"common"}, "404", "Not found")
		r.JSON(http.StatusNotFound, FormatErr(errors))
	} else {
		r.HTML(http.StatusNotFound, "error/404", nil)
	}
}

func APIEntry(r render.Render) {
	r.JSON(http.StatusOK, map[string]interface{}{
		"tokens":  "/api/v1/tokens",
		"users":   "/api/v1/users",
		"domains": "/api/v1/domains",
		"records": "/api/v1/records",
		"emails":  "/api/v1/emails",
	})
}
