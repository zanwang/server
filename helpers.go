package main

import (
  "html/template"
)

var AppHelpers = template.FuncMap{
  "isset": func(data map[string]interface{}, name string) bool {
    _, ok := data[name]
    return ok
  },
  "mapValue": func(data map[string]interface{}, name string) interface{} {
    return data[name]
  },
}