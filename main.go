package main

import (
	c "account/http/container"
	"account/http/server"
	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	container := c.NewContainer()

	defer container.AMQPService.Close()

	engine, err := server.SetupServer(container)
	if err != nil {
		log.Fatal(err.Error())
	}
	port := os.Getenv("APP_PORT")

	serverErr := engine.Run(":" + port)
	if serverErr != nil {
		log.Fatal(serverErr.Error())
	}
}
