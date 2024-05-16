package server

import (
	c "account/http/container"
	"account/http/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
)

func SetupTestServer(wg *sync.WaitGroup) {
	engine := gin.Default()

	container := c.NewContainer()

	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, "ok")
	})

	routes.RegisterUserRoutes(engine, container)
	routes.RegisterAppRoutes(engine, container)
	routes.RegisterAuthRoutes(engine, container)
	routes.RegisterTestRoutes(engine, container)
	routes.RegisterSwaggerRoutes(engine)

	port := os.Getenv("APP_PORT")

	go func() {
		if err := engine.Run(":" + port); err != nil {
			fmt.Println("Failed to start the server:", err)
			return
		}
	}()

	wg.Done()
}
