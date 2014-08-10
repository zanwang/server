package middlewares

import (
	"net/http"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/tommy351/maji.moe/controllers"
)

func Validate(obj interface{}, ifacePtr ...interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		contentType := req.Header.Get("Content-Type")
		if req.Method == "POST" || req.Method == "PUT" || contentType != "" {
			if strings.Contains(contentType, "form-urlencoded") {
				context.Invoke(binding.Form(obj, ifacePtr...))
			} else if strings.Contains(contentType, "multipart/form-data") {
				context.Invoke(binding.MultipartForm(obj, ifacePtr...))
			} else if strings.Contains(contentType, "json") {
				context.Invoke(binding.Json(obj, ifacePtr...))
			} else {
				var errors binding.Errors
				if contentType == "" {
					errors.Add([]string{}, binding.ContentTypeError, "Empty Content-Type")
				} else {
					errors.Add([]string{}, binding.ContentTypeError, "Unsupported Content-Type")
				}
				context.Map(errors)
			}
		} else {
			context.Invoke(binding.Form(obj, ifacePtr...))
		}

		context.Invoke(FormErrorHandler)
	}
}

func FormErrorHandler(errors binding.Errors, r render.Render) {
	if errors != nil {
		r.JSON(http.StatusBadRequest, controllers.FormatErr(errors))
		return
	}
}
