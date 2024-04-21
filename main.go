package main

import (
	"account/http/server"
	"log"
	"os"
)

// @title Account API
// @version 0.5.0
// @contact.email prodemmi@gmail.com
// @BasePath /api/v1
// @accept json
// @produce json
// @schemes http https
func main() {
	engine, err := server.SetupServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	port := os.Getenv("APP_PORT")

	serverErr := engine.Run(":" + port)
	if serverErr != nil {
		log.Fatal(serverErr.Error())
	}
}
