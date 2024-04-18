package routes

import (
	"account/http/middlewares"
	"account/internal"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(app *gin.Engine, container *internal.Container) {
	appMiddleware := middlewares.NewAppMiddleware(container.AppService)
	jwtMiddleware := middlewares.NewJWTMiddleware(container.UserService)

	v1 := app.Group("/v1/auth/").Use(appMiddleware)
	{
		v1.GET("/me", jwtMiddleware, container.AuthController.Me)

		v1.POST("/sign-up", container.AuthController.SignUp)
		v1.PUT("/sign-in", container.AuthController.SignIn)

		v1.POST("/verify-otp", container.AuthController.VerifyOtp)
	}
}
