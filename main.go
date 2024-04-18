package main

import (
	"account/internal"
	"account/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// @title Account API

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	port := os.Getenv("APP_PORT")

	serverErr := engine.Run(":" + port)
	if serverErr != nil {
		fmt.Print(serverErr.Error())
		return
	}

}
