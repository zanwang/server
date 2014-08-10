package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/method"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/controllers"
	"github.com/tommy351/maji.moe/middlewares"
	"github.com/tommy351/maji.moe/models"
)

func main() {
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

	// Routes
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/app", controllers.App)
	m.Get("/login", controllers.App)
	m.Get("/signup", controllers.App)

	m.Group("/api/v1", func(r martini.Router) {
		r.Group("/tokens", func(r martini.Router) {
			// POST /api/v1/tokens
			r.Post("", middlewares.Validate(controllers.TokenCreateForm{}), controllers.TokenCreate)
			// DELETE /api/v1/tokens
			r.Delete("", middlewares.NeedAuthorization, controllers.TokenDestroy)
		}, middlewares.CheckToken)

		r.Group("/users", func(r martini.Router) {
			// POST /api/v1/users
			r.Post("", middlewares.Validate(controllers.UserCreateForm{}), controllers.UserCreate)

			r.Group("/:user_id", func(r martini.Router) {
				// GET /api/v1/users/:user_id
				r.Get("", controllers.UserShow)

				r.Group("", func(r martini.Router) {
					// PUT /api/v1/users/:user_id
					r.Put("", middlewares.Validate(controllers.UserUpdateForm{}), controllers.UserUpdate)
					// DELETE /api/v1/users/:user_id
					r.Delete("", controllers.UserDestroy)
					// GET /api/v1/users/:user_id/domains
					r.Get("/domains", controllers.DomainList)
					// POST /api/v1/users/:user_id/domains
					r.Post("/domains", controllers.DomainCreate)
				}, middlewares.NeedAuthorization)
			}, middlewares.CheckToken)
		})

		r.Group("/domains", func(r martini.Router) {
			r.Group("/:domain_id", func(r martini.Router) {
				// GET /api/v1/domains/:domain_id
				r.Get("", controllers.DomainShow)
				// PUT /api/v1/domains/:domain_id
				r.Put("", controllers.DomainUpdate)
				// DELETE /api/v1/domains/:domain_id
				r.Delete("", controllers.DomainDestroy)
				// GET /api/v1/domain/:domain_id/records
				r.Get("/records", controllers.RecordList)
				// POST /api/v1/domain/:domain_id/records
				r.Post("/records", controllers.RecordCreate)
			})
		}, middlewares.CheckToken, middlewares.NeedAuthorization)

		r.Group("/records", func(r martini.Router) {
			r.Group("/:record_id", func(r martini.Router) {
				// GET /api/v1/records/:record_id
				r.Get("", controllers.RecordShow)
				// PUT /api/v1/records/:record_id
				r.Put("", controllers.RecordUpdate)
				// DELETE /api/v1/records/:record_id
				r.Delete("", controllers.RecordDestroy)
			})
		}, middlewares.CheckToken, middlewares.NeedAuthorization)
	})

	// Serve static files
	m.Use(martini.Static("public"))

	// Start server
	log.Printf("Listening on http://%s:%d", host, port)
	panic(http.ListenAndServe(host+":"+strconv.Itoa(port), m))
}
