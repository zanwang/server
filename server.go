package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tommy351/maji.moe/auth"
	"github.com/tommy351/maji.moe/config"
	"github.com/tommy351/maji.moe/controllers"
	"github.com/tommy351/maji.moe/middlewares"
	"github.com/tommy351/maji.moe/models"
)

func server() {
	// Load configuration
	config := config.Load()

	// Load database
	dbMap := models.Load(config)
	defer dbMap.Db.Close()

	// Load mailgun
	mg := mail(config)

	// OAuth
	fbApp := auth.LoadFacebook(config)
	twitterConsumer := auth.LoadTwitter(config)

	// Create a classic martini
	m := martini.Classic()
	host := config.Server.Host
	port := config.Server.Port
	secret := config.Server.Secret

	// Basic setup
	m.Map(config)
	m.Map(dbMap)
	m.Map(mg)
	m.Map(fbApp)
	m.Map(twitterConsumer)

	// Middlewares
	store := sessions.NewCookieStore([]byte(secret))

	m.Use(sessions.Sessions("my_session", store))
	m.Use(render.Renderer(render.Options{
		Directory:  "views",
		Extensions: []string{".html", ".htm"},
		Funcs:      []template.FuncMap{appHelpers},
	}))

	// Routes
	m.Get("/", controllers.Home)
	m.Get("/app", controllers.App)
	m.Get("/login", controllers.App)
	m.Get("/signup", controllers.App)
	m.Get("/forgot_password", controllers.App)
	m.Get("/settings", controllers.App)
	m.Get("/activation/:token", controllers.UserActivate)
	m.NotFound(controllers.NotFound)

	m.Group("/api/v1", func(r martini.Router) {
		r.Get("", controllers.APIEntry)

		r.Group("/tokens", func(r martini.Router) {
			// POST /api/v1/tokens
			r.Post("", middlewares.Validate(controllers.TokenCreateForm{}), controllers.TokenCreate, middlewares.ResponseToken)
			// PUT /api/v1/tokens
			r.Put("", middlewares.CheckToken, middlewares.NeedAuthorization, controllers.TokenUpdate, middlewares.ResponseToken)
			// DELETE /api/v1/tokens
			r.Delete("", middlewares.CheckToken, middlewares.NeedAuthorization, controllers.TokenDestroy)
			// POST /api/v1/tokens/facebook
			r.Post("/facebook", middlewares.Validate(controllers.TokenFacebookForm{}), controllers.TokenFacebook, middlewares.ResponseToken)
		})

		r.Group("/users", func(r martini.Router) {
			// POST /api/v1/users
			r.Post("", middlewares.Validate(controllers.UserCreateForm{}), controllers.UserCreate)

			r.Group("/:user_id", func(r martini.Router) {
				// GET /api/v1/users/:user_id
				r.Get("", controllers.UserShow)

				r.Group("", func(r martini.Router) {
					// PUT /api/v1/users/:user_id
					r.Put("", middlewares.CheckCurrentUser, middlewares.Validate(controllers.UserUpdateForm{}), controllers.UserUpdate)
					// DELETE /api/v1/users/:user_id
					r.Delete("", middlewares.CheckCurrentUser, controllers.UserDestroy)
					// GET /api/v1/users/:user_id/domains
					r.Get("/domains", controllers.DomainList)
					// POST /api/v1/users/:user_id/domains
					r.Post("/domains", middlewares.CheckCurrentUser, middlewares.NeedActivation, middlewares.Validate(controllers.DomainCreateForm{}), controllers.DomainCreate)
				})
			}, middlewares.CheckToken, middlewares.NeedAuthorization, middlewares.GetUser)
		})

		r.Group("/domains", func(r martini.Router) {
			r.Group("/:domain_id", func(r martini.Router) {
				// GET /api/v1/domains/:domain_id
				r.Get("", middlewares.CheckOwnershipOfDomain, controllers.DomainShow)
				// PUT /api/v1/domains/:domain_id
				r.Put("", middlewares.CheckOwnershipOfDomain, middlewares.Validate(controllers.DomainUpdateForm{}), controllers.DomainUpdate)
				// DELETE /api/v1/domains/:domain_id
				r.Delete("", middlewares.CheckOwnershipOfDomain, controllers.DomainDestroy)
				// GET /api/v1/domain/:domain_id/records
				r.Get("/records", middlewares.CheckOwnershipOfDomain, controllers.RecordList)
				// POST /api/v1/domain/:domain_id/records
				r.Post("/records", middlewares.CheckOwnershipOfDomain, middlewares.Validate(controllers.RecordCreateForm{}), controllers.RecordCreate)
			}, middlewares.GetDomain)
		}, middlewares.CheckToken, middlewares.NeedAuthorization)

		r.Group("/records", func(r martini.Router) {
			r.Group("/:record_id", func(r martini.Router) {
				// GET /api/v1/records/:record_id
				r.Get("", middlewares.CheckOwnershipOfRecord, controllers.RecordShow)
				// PUT /api/v1/records/:record_id
				r.Put("", middlewares.CheckOwnershipOfRecord, middlewares.Validate(controllers.RecordUpdateForm{}), controllers.RecordUpdate)
				// DELETE /api/v1/records/:record_id
				r.Delete("", middlewares.CheckOwnershipOfRecord, controllers.RecordDestroy)
			}, middlewares.GetRecord)
		}, middlewares.CheckToken, middlewares.NeedAuthorization)

		r.Group("/emails", func(r martini.Router) {
			r.Post("/resend", middlewares.Validate(controllers.EmailResendForm{}), controllers.EmailResend)
		})
	})

	m.Group("/oauth", func(r martini.Router) {
		m.Group("/twitter", func(r martini.Router) {
			// GET /oauth/twitter/login
			m.Get("/login", controllers.OAuthTwitterLogin)
			// GET /oauth/twitter/callback
			m.Get("/callback", controllers.OAuthTwitterCallback)
		})

		m.Group("/google", func(r martini.Router) {
			//
		})
	})

	// Start server
	log.Printf("Listening on http://%s:%d", host, port)
	panic(http.ListenAndServe(host+":"+strconv.Itoa(port), m))
}

func main() {
	server()
}
