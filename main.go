package main

import (
	"account/http"
	"account/internal"
	"log"
	"os"
)

// @title Account API
func main() {
	container := internal.NewContainer()

	engine, err := http.SetupServer(container)
	if err != nil {
		log.Fatal(err.Error())
	}
	port := os.Getenv("APP_PORT")

	serverErr := engine.Run(":" + port)
	if serverErr != nil {
		log.Fatal(serverErr.Error())
	}
}
