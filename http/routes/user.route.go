package routes

import (
	"account/http/container"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(app *gin.Engine, container *container.Container) {

	v1 := app.Group("/v1/users/").Use(container.AppMiddleware)
	{
		v1.POST("/", container.UserController.Store)
		v1.GET("/", container.UserController.Index)
		v1.GET("/:id", container.UserController.Single)
	}
}
