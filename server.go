package main

import (
  "net/http"
  "log"
  "strconv"
  "html/template"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/sessions"
  "github.com/martini-contrib/method"
  "github.com/martini-contrib/binding"
  "github.com/coopernurse/gorp"
)

var DbMap *gorp.DbMap

func main() {
  m := martini.Classic()
  config := Config.Get("server").(map[interface{}]interface{})
  host := config["host"].(string)
  port := config["port"].(int)
  secret := config["secret"].(string)

  // Initialize database
  DbMap = initDb()
  defer DbMap.Db.Close()

  // Middlewares
  m.Use(render.Renderer(render.Options{
    Directory: "views",
    Extensions: []string{".html", ".htm"},
    Funcs: []template.FuncMap{AppHelpers},
  }))

  store := sessions.NewCookieStore([]byte(secret))
  m.Use(sessions.Sessions("my_session", store))

  m.Use(method.Override())

  // Routes
  m.Get("/", func(r render.Render) {
    r.HTML(200, "index", nil)
  })

  m.Get("/login", GetCurrentUser, SessionNew)
  m.Post("/login", binding.Form(SessionCreateForm{}), SessionCreate)
  m.Get("/logout", SessionDestroy)

  m.Get("/signup", GetCurrentUser, UserNew)
  m.Post("/signup", binding.Form(UserCreateForm{}), UserCreate)
  m.Get("/users/:user_id", GetCurrentUser, UserShow)
  m.Put("/users/:user_id", GetCurrentUser, UserUpdate)
  m.Delete("/users/:user_id", GetCurrentUser, UserDestroy)
  m.Get("/users/:user_id/confirm", UserConfirm)

  m.Get("/settings/profile", GetCurrentUser, UserEdit)

  m.Post("/domains", DomainCreate)
  m.Get("/domains/:domain_id", DomainShow)
  m.Get("/domains/:domain_id/edit", DomainEdit)
  m.Put("/domains/:domain_id", DomainUpdate)
  m.Delete("/domains/:domain_id", DomainDestroy)

  m.Post("/domains/:domain_id/records", RecordCreate)
  m.Put("/domains/:domain_id/records/:record_id", RecordUpdate)
  m.Delete("/domains/:domain_id/records/:record_id", RecordDestroy)

  // Serve static files
  m.Use(martini.Static("public"))

  log.Printf("Listening at http://%s:%d", host, port)
  log.Fatal(http.ListenAndServe(host + ":" + strconv.Itoa(port), m))
}

func checkErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}

func formatErr(errors binding.Errors) map[string]interface{} {
  result := make(map[string]interface{})

  for _, err := range errors {
    for _, field := range err.Fields() {
      if _, ok := result[field]; !ok {
        result[field] = err.Error()
      }
    }
  }

  return result
}