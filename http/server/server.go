package server

import (
	c "account/http/container"
	"account/http/routes"
	"account/internal"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupServer() (*gin.Engine, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	engine := gin.Default()

	container := c.NewContainer()

	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, "ok")
	})

	engine.GET("/version", func(context *gin.Context) {
		context.JSON(http.StatusOK, "0.6.2")
	})

	engine.Use(cors.Default())

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
