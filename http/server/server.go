package server

import (
	"account/http/container"
	"account/http/routes"
	"account/internal"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func SetupServer(container *container.Container) (*gin.Engine, error) {
	engine := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "APP-ID", "API-TOKEN", "Content-Type", "Access-Control-Allow-Origin", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	engine.Use(cors.New(config))

	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, "ok")
	})

	engine.GET("/version", func(context *gin.Context) {
		context.JSON(http.StatusOK, "0.7.0")
	})

	routes.RegisterUserRoutes(engine, container)
	routes.RegisterAppRoutes(engine, container)
	routes.RegisterAuthRoutes(engine, container)
	routes.RegisterSwaggerRoutes(engine)

	dbErr := internal.SetupDatabase()
	if dbErr != nil {
		return nil, dbErr
	}

	container.AppAMQP.Setup()

	return engine, nil
}
