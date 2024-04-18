package main

import (
	"account/http/server"
	"log"
	"os"
)

// @title Account API
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
