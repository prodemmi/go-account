package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// @title Account API
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	server, err := SetupServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	port := os.Getenv("APP_PORT")

	serverErr := server.Run(":" + port)
	if serverErr != nil {
		log.Fatal(serverErr.Error())
	}
}
