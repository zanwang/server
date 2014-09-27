package server

import (
	"path"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/config"
	"github.com/majimoe/server/controllers"
	"github.com/tommy351/gin-cors"
	"github.com/tommy351/gin-csrf"
	"github.com/tommy351/gin-sessions"
)

func Server() *gin.Engine {
	conf := config.Config
	r := gin.New()
	store := sessions.NewCookieStore([]byte(conf.Server.Secret))
	middleware := controllers.Middleware{}

	if conf.Server.Logger {
		r.Use(gin.Logger())
	}

	r.Use(controllers.Recovery)

	r.Use(sessions.Middleware("mm_session", store))

	r.Use(csrf.Middleware(csrf.Options{
		Secret: conf.Server.Secret,
	}))

	r.Use(cors.Middleware(cors.Options{
		MaxAge:       time.Hour * 24,
		AllowHeaders: []string{"Origin", "Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
	}))

	r.GET("/", controllers.Home)
	r.GET("/app", controllers.App)
	r.GET("/app/settings", controllers.App)
	r.GET("/app/domains/:id", controllers.App)
	r.GET("/login", controllers.App)
	r.GET("/signup", controllers.App)
	r.GET("/password-reset", controllers.App)
	r.GET("/users/:user_id/activation/:token", controllers.UserActivation)
	r.GET("/users/:user_id/passwords/reset/:token", controllers.PasswordReset)
	r.POST("/users/:user_id/passwords/reset/:token", controllers.PasswordResetSubmit)
	r.NotFound404(controllers.NotFound)

	apiv1Group := r.Group("/api/v1")
	{
		apiv1 := controllers.APIv1{}

		apiv1Group.Use(apiv1.Recovery)
		apiv1Group.Use(apiv1.UpdateToken)

		apiv1Group.GET("", apiv1.Entry)
		apiv1Group.POST("/tokens", apiv1.TokenCreate)
		apiv1Group.PUT("/tokens", middleware.TokenRequired, apiv1.TokenUpdate)
		apiv1Group.DELETE("/tokens", middleware.TokenRequired, apiv1.TokenDestroy)

		apiv1Group.POST("/users", apiv1.UserCreate)
		apiv1Group.GET("/users/:user_id", middleware.GetUser, middleware.GetToken, apiv1.UserShow)
		apiv1Group.PUT("/users/:user_id", middleware.TokenRequired, middleware.CheckPermissionOfUser, apiv1.UserUpdate)
		apiv1Group.DELETE("/users/:user_id", middleware.TokenRequired, middleware.CheckPermissionOfUser, apiv1.UserDestroy)
		apiv1Group.GET("/users/:user_id/domains", middleware.GetUser, apiv1.DomainList)
		apiv1Group.POST("/users/:user_id/domains", middleware.TokenRequired, middleware.CheckPermissionOfUser, apiv1.DomainCreate)
		apiv1Group.GET("/users/:user_id/tokens", middleware.TokenRequired, middleware.CheckPermissionOfUser, apiv1.TokenList)

		apiv1Group.GET("/domains", apiv1.DomainList)
		apiv1Group.GET("/domains/:domain_id", middleware.GetDomain, apiv1.DomainShow)
		apiv1Group.PUT("/domains/:domain_id", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.DomainUpdate)
		apiv1Group.DELETE("/domains/:domain_id", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.DomainDestroy)
		apiv1Group.POST("/domains/:domain_id/renew", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.DomainRenew)
		apiv1Group.GET("/domains/:domain_id/records", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.RecordList)
		apiv1Group.POST("/domains/:domain_id/records", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.RecordCreate)

		apiv1Group.GET("/records/:record_id", middleware.TokenRequired, middleware.CheckOwnershipOfRecord, apiv1.RecordShow)
		apiv1Group.PUT("/records/:record_id", middleware.TokenRequired, middleware.CheckOwnershipOfRecord, apiv1.RecordUpdate)
		apiv1Group.DELETE("/records/:record_id", middleware.TokenRequired, middleware.CheckOwnershipOfRecord, apiv1.RecordDestroy)

		apiv1Group.POST("/emails/resend", apiv1.EmailResend)
		apiv1Group.POST("/passwords/reset", apiv1.PasswordReset)
	}

	r.Use(static.Serve(path.Join(config.BaseDir, "public")))

	return r
}
