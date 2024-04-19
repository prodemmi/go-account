package routes

import (
	"account/http/container"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(app *gin.Engine, container *container.Container) {
	v1 := app.Group("/v1/auth/").Use(container.AppMiddleware)
	{
		v1.GET("/me", container.JWTMiddleware, container.AuthController.Me)

		v1.POST("/sign-up", container.AuthController.SignUp)
		v1.PUT("/sign-in", container.AuthController.SignIn)

		v1.POST("/verify-otp", container.AuthController.VerifyOtp)
	}
}
