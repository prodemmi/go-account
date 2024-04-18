package server

import (
	c "account/http/container"
	"account/http/routes"
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

	container := c.NewContainer()

	routes.RegisterUserRoutes(engine, container)
	routes.RegisterAppRoutes(engine, container)
	routes.RegisterAuthRoutes(engine, container)
	routes.RegisterSwaggerRoutes(engine)

	return engine, nil
}
