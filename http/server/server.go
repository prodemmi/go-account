package server

import (
	c "account/http/container"
	"account/http/routes"
	"account/internal"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupServer() (*gin.Engine, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	engine := gin.Default()

	container := c.NewContainer()

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

	return engine, nil
}
