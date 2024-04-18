package main

import (
	"account/internal"
	"account/routes"
	"github.com/gin-gonic/gin"
)

func SetupServer() (*gin.Engine, error) {
	engine := gin.Default()

	routes.RegisterUserRoutes(engine)
	routes.RegisterAppRoutes(engine)
	routes.RegisterAuthRoutes(engine)
	routes.RegisterSwaggerRoutes(engine)

	dbErr := internal.SetupDatabase()
	if dbErr != nil {
		return nil, dbErr
	}

	return engine, nil
}
