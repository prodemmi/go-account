package routes

import (
	"account/http/container"
	"account/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(app *gin.Engine, container *container.Container) {

	v1 := app.Group("/v1/users/").Use(middlewares.NewAppMiddleware(container.AppService))
	{
		v1.GET("/", container.UserController.Index)
		v1.GET("/:id", container.UserController.Single)

		v1.POST("/", container.UserController.Store)
	}
}
