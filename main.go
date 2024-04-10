package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"account/internal"
	"account/routes"
)

func main() {
	engine := gin.Default()
	db, dbErr := internal.SetupDatabase()
	if dbErr != nil {
		fmt.Print(dbErr.Error())
		return
	}

	routes.RegisterUserRoutes(engine, db)
	routes.RegisterAppRoutes(engine, db)

	serverErr := engine.Run(":8081")
	if serverErr != nil {
		fmt.Print(serverErr.Error())
		return
	}

}
