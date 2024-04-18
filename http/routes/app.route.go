package routes

import (
	"account/http/container"
	"github.com/gin-gonic/gin"
)

func RegisterAppRoutes(app *gin.Engine, container *container.Container) {
	v1 := app.Group("/v1/apps/")
	{
		v1.POST("/", container.AppController.Store)
		v1.GET("/", container.AppController.Index)
		v1.GET("/:id", container.AppController.Single)
		v1.DELETE("/:id", container.AppController.Delete)
	}
}
