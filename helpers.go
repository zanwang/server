package main

import (
	"html/template"
)

var appHelpers = template.FuncMap{
	"isset": func(data map[string]interface{}, name string) bool {
		_, ok := data[name]
		return ok
	},
	"map": func(data map[string]interface{}, name string) interface{} {
		return data[name]
	},
}
