package http

import (
	"account/http/routes"
	"account/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"path/filepath"
)

func SetupTestServer() (*gin.Engine, error) {
	absPath, _ := filepath.Abs("../")

	if err := godotenv.Load(fmt.Sprintf("%v/.env.test", absPath)); err != nil {
		log.Fatal(err.Error())
	}

	engine := gin.Default()

	container := internal.NewContainer()

	routes.RegisterUserRoutes(engine, container)
	routes.RegisterAppRoutes(engine, container)
	routes.RegisterAuthRoutes(engine, container)
	routes.RegisterSwaggerRoutes(engine)

	return engine, nil
}
