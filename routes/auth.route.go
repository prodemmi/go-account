package routes

import (
	"account/controllers"
	"account/middlewares"
	"account/repositories"
	"account/services"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(app *gin.Engine) {
	authRepository := repositories.NewAuthRepository()
	autService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(autService)

	appMiddleware := middlewares.NewAppMiddleware()
	jwtMiddleware := middlewares.NewJWTMiddleware()

	v1 := app.Group("/v1/auth/").Use(appMiddleware)
	{
		v1.GET("/me", jwtMiddleware, authController.Me)
	}
}
