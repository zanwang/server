package controllers

import (
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"
  "encoding/xml"
  "github.com/tommy351/maji.moe/util"
)

type Controller struct {
  w http.ResponseWriter
  r *http.Request
  methods map[string]func()
}

type ControllerInterface interface {
  Handle(w http.ResponseWriter, r *http.Request)
  Get()
  Post()
  Put()
  Delete()
  Patch()
  Head()
  Error(err error)
  Init()
  AddMethod(name string, fn func())
}

func (c *Controller) Init() {
  c.methods = make(map[string]func())
}

func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) {
  c.w = w
  c.r = r

  if fn, ok := c.methods[r.Method]; ok {
    fn()
  } else {
    c.methodNotAllowed()
  }
}

func (c *Controller) Get() {
  c.methodNotAllowed()
}

func (c *Controller) Post() {
  c.methodNotAllowed()
}

func (c *Controller) Put() {
  c.methodNotAllowed()
}

func (c *Controller) Delete() {
  c.methodNotAllowed()
}

func (c *Controller) Patch() {
  c.methodNotAllowed()
}

func (c *Controller) Head() {
  c.methodNotAllowed()
}

func (c *Controller) methodNotAllowed() {
  http.Error(c.w, "Method Not Allowed", 405)
}

func (c *Controller) Write(data interface{}) {
  fmt.Fprint(c.w, data)
}

func (c *Controller) Render(path string, data interface{}) {
  tmpl, err := template.ParseFiles(util.ResolveView(path))

  if err != nil {
    c.Error(err)
    return
  }

  c.Template(tmpl, data)
}

func (c *Controller) Template(tmpl *template.Template, data interface{}) {
  if err := tmpl.Execute(c.w, data); err != nil {
    c.Error(err)
  }
}

func (c *Controller) Json(data interface{}) {
  var (
    b []byte
    err error
  )

  env := util.Environment()

  if env == "prod" {
    b, err = json.Marshal(data)
  } else {
    b, err = json.MarshalIndent(data, "", "  ")
  }

  if err != nil {
    c.Error(err)
    return
  }

  c.Write(b)
}

func (c *Controller) Xml(data interface{}) {
  var (
    b []byte
    err error
  )

  env := util.Environment()

  if env == "prod" {
    b, err = xml.Marshal(data)
  } else {
    b, err = xml.MarshalIndent(data, "", "  ")
  }

  if err != nil {
    c.Error(err)
    return
  }

  c.Write(b)
}

func (c *Controller) Redirect(path string, code int) {
  http.Redirect(c.w, c.r, path, code)
}

func (c *Controller) Error(err error) {
  http.Error(c.w, err.Error(), http.StatusInternalServerError)
}

func (c *Controller) AddMethod(name string, fn func()) {
  c.methods[name] = fn
}