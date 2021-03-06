package util

import (
	"path"

	"github.com/majimoe/server/config"
	"gopkg.in/unrolled/render.v1"
)

var Render *render.Render

func init() {
	Render = render.New(render.Options{
		Directory:     path.Join(config.BaseDir, "views"),
		Extensions:    []string{".html", ".htm"},
		IsDevelopment: config.Env == config.Development,
	})
}
