package routes

import (
	"account/http/middlewares"
	"account/internal"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(app *gin.Engine, container *internal.Container) {

	v1 := app.Group("/v1/users/").Use(middlewares.NewAppMiddleware(container.AppService))
	{
		v1.GET("/", container.UserController.Index)
		v1.GET("/:id", container.UserController.Single)

		v1.POST("/", container.UserController.Store)
	}
}
