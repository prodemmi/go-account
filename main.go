package main

import (
	"account/internal"
	"account/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

// @title Account API

func main() {
	engine := gin.Default()
	dbErr := internal.SetupDatabase()
	if dbErr != nil {
		fmt.Print(dbErr.Error())
		return
	}

	routes.RegisterUserRoutes(engine)
	routes.RegisterAppRoutes(engine)
	routes.RegisterAuthRoutes(engine)
	routes.RegisterSwaggerRoutes(engine)

	serverErr := engine.Run(":8081")
	if serverErr != nil {
		fmt.Print(serverErr.Error())
		return
	}

}
