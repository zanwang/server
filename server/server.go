package server

import (
	"path"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/config"
	"github.com/majimoe/server/controllers"
)

func Server() *gin.Engine {
	conf := config.Config
	r := gin.New()

	if conf.Server.Logger {
		r.Use(gin.Logger())
	}

	r.Use(controllers.Recovery)

	r.GET("/", controllers.Home)
	r.GET("/app", controllers.App)
	r.GET("/login", controllers.App)
	r.GET("/signup", controllers.App)
	r.GET("/forgot_password", controllers.App)
	r.GET("/settings", controllers.App)
	r.GET("/users/:user_id/activation/:token", controllers.UserActivation)
	r.NotFound404(controllers.NotFound)

	apiv1Group := r.Group("/api/v1")
	middleware := controllers.Middleware{}

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

		apiv1Group.GET("/domains/:domain_id", middleware.GetDomain, apiv1.DomainShow)
		apiv1Group.PUT("/domains/:domain_id", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.DomainUpdate)
		apiv1Group.DELETE("/domains/:domain_id", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.DomainDestroy)
		apiv1Group.POST("/domains/:domain_id/renew", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.DomainRenew)
		apiv1Group.GET("/domains/:domain_id/records", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.RecordList)
		apiv1Group.POST("/domains/:domain_id/records", middleware.TokenRequired, middleware.CheckOwnershipOfDomain, apiv1.RecordCreate)

		apiv1Group.GET("/records/:record_id", middleware.TokenRequired, middleware.CheckOwnershipOfRecord, apiv1.RecordShow)
		apiv1Group.PUT("/records/:record_id", middleware.TokenRequired, middleware.CheckOwnershipOfRecord, apiv1.RecordUpdate)
		apiv1Group.DELETE("/records/:record_id", middleware.TokenRequired, middleware.CheckOwnershipOfRecord, apiv1.RecordDestroy)
	}

	r.Use(static.Serve(path.Join(config.BaseDir, "public")))

	return r
}
