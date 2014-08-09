package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/method"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/controllers"
	"github.com/tommy351/maji.moe/middlewares"
	"github.com/tommy351/maji.moe/models"
)

func main() {
	// Recovery
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Load configuration
	config := config.Load()

	// Load database
	dbMap := models.Load(config)
	defer dbMap.Db.Close()

	// Createa a classic martini
	m := martini.Classic()
	host := config.Server.Host
	port := config.Server.Port
	secret := config.Server.Secret

	// Mapping
	m.Map(config)
	m.Map(dbMap)

	// Middlewares
	store := sessions.NewCookieStore([]byte(secret))

	m.Use(sessions.Sessions("my_session", store))
	m.Use(render.Renderer(render.Options{
		Directory:  "views",
		Extensions: []string{".html", ".htm"},
		Funcs:      []template.FuncMap{appHelpers},
	}))
	m.Use(method.Override())
	m.Use(middlewares.CSRFMiddleware)

	// Routes
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/app", func(r render.Render) {
		r.HTML(200, "app", nil)
	})

	m.Get("/login", middlewares.GetCurrentUser, controllers.SessionNew)
	m.Post("/sessions", binding.Form(controllers.SessionCreateForm{}), controllers.SessionCreate)
	m.Get("/logout", middlewares.GetCurrentUser, controllers.SessionDestroy)
	m.Delete("/sessions", middlewares.GetCurrentUser, controllers.SessionDestroy)

	m.Get("/signup", middlewares.GetCurrentUser, controllers.UserNew)
	m.Post("/users", binding.Form(controllers.UserCreateForm{}), controllers.UserCreate)

	m.Group("/users/", func(r martini.Router) {
		r.Put("/:id", controllers.UserUpdate)
		r.Delete("/:id", controllers.UserDestroy)
	}, middlewares.GetCurrentUser, middlewares.NeedLogin)

	m.Get("/accounts", middlewares.GetCurrentUser, middlewares.NeedLogin, controllers.UserEdit)

	// Serve static files
	m.Use(martini.Static("public"))

	// Start server
	log.Printf("Listening on http://%s:%d", host, port)
	panic(http.ListenAndServe(host+":"+strconv.Itoa(port), m))
}
