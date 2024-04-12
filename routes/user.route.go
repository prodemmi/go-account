package routes

import (
	"account/controllers"
	"account/middlewares"
	"account/repositories"
	"account/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(app *gin.Engine) {
	profileRepository := repositories.NewProfileRepository()
	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(userRepository)
	profileService := services.NewProfileService(profileRepository)
	userController := controllers.NewUserController(userService, profileService)

	v1 := app.Group("/v1/users/").Use(middlewares.NewAppMiddleware())
	{
		v1.GET("/", userController.Index)
		v1.GET("/:id", userController.Single)

		v1.POST("/", userController.Store)
	}
}
