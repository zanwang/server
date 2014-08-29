package util

import (
	"github.com/tommy351/maji.moe/config"
	"gopkg.in/unrolled/render.v1"
)

var Render *render.Render

func init() {
	Render = render.New(render.Options{
		Directory:     "views",
		Extensions:    []string{".html", ".htm"},
		IsDevelopment: config.Env == config.Development,
	})
}
