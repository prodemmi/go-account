package http

import (
	"account/http/routes"
	"account/internal"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func SetupServer() (*gin.Engine, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	engine := gin.Default()

	container := internal.NewContainer()

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
