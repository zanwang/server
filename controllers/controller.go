package controllers

import (
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"
  "encoding/xml"
  "log"
  "reflect"
  "github.com/tommy351/maji.moe/util"
)

type Controller struct {
  w http.ResponseWriter
  r *http.Request
  methods map[string]func(w http.ResponseWriter, r *http.Request)
  filters struct {
    before map[string][]interface{}
    after map[string][]interface{}
  }
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
  Prepare()
  Init()
  AddMethod(name string, fn func())
}

func (c *Controller) Prepare() {
  c.methods = make(map[string]func(w http.ResponseWriter, r *http.Request))
}

func (c *Controller) Init() {
  // Register filters here
}

func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) {
  c.w = w
  c.r = r

  if fn, ok := c.methods[r.Method]; ok {
    fn(w, r)
  } else {
    http.NotFound(c.w, c.r)
  }
}

func (c *Controller) Before(method string, fn ...interface{}) {
  c.filters.before[method] = fn
}

func (c *Controller) After(method string, fn ...interface{}) {
  c.filters.after[method] = fn
}

func noopHandler(w http.ResponseWriter, r *http.Request){}

// https://github.com/shelakel/go-middleware
func (c *Controller) AddMethod(method string, fn func()) {
  filters := []interface{}{fn}
  handler := noopHandler

  if before, ok := c.filters.before[method]; ok {
    filters = append(before, filters)
  }

  if after, ok := c.filters.after[method]; ok {
    filters = append(filters, after)
  }

  for i := len(filters) - 1; i >= 0; i-- {
    next := handler

    switch current := filters[i].(type) {
    case func(http.ResponseWriter, *http.Request, func()):
      handler = func(w http.ResponseWriter, r *http.Request) {
        current(w, r, func() {
          next(w, r)
        })
      }
    case func(http.ResponseWriter, *http.Request):
      handler = func(w http.ResponseWriter, r *http.Request) {
        current(w, r)
        next(w, r)
      }
    case func(*Controller, func()):
      handler = func(w http.ResponseWriter, r *http.Request) {
        current(c, func() {
          next(w, r)
        })
      }
    case func(*Controller):
      handler = func(w http.ResponseWriter, r *http.Request) {
        current(c)
        next(w, r)
      }
    case func(func()):
      handler = func(w http.ResponseWriter, r *http.Request) {
        current(func() {
          next(w, r)
        })
      }
    case func():
      handler = func(w http.ResponseWriter, r *http.Request) {
        current()
        next(w, r)
      }
    default:
        log.Panicf("Unsupported middleware type '%v' at index %d", reflect.TypeOf(current), i)
    }
  }

  c.methods[method] = handler
}

func (c *Controller) Get() {
  http.NotFound(c.w, c.r)
}

func (c *Controller) Post() {
  http.NotFound(c.w, c.r)
}

func (c *Controller) Put() {
  http.NotFound(c.w, c.r)
}

func (c *Controller) Delete() {
  http.NotFound(c.w, c.r)
}

func (c *Controller) Patch() {
  http.NotFound(c.w, c.r)
}

func (c *Controller) Head() {
  http.NotFound(c.w, c.r)
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