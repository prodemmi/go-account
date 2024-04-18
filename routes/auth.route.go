package routes

import (
	"account/controllers"
	"account/middlewares"
	"account/repositories"
	"account/services"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(app *gin.Engine) {
	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(userRepository)
	authRepository := repositories.NewAuthRepository()
	autService := services.NewAuthService(authRepository, userService)
	otpRepository := repositories.NewOTPRepository()
	otpService := services.NewOTPService(otpRepository)
	authController := controllers.NewAuthController(autService, userService, otpService)

	appMiddleware := middlewares.NewAppMiddleware()
	jwtMiddleware := middlewares.NewJWTMiddleware()

	v1 := app.Group("/v1/auth/").Use(appMiddleware)
	{
		v1.GET("/me", jwtMiddleware, authController.Me)

		v1.POST("/sign-up", authController.SignUp)
		v1.PUT("/sign-in", authController.SignIn)

		v1.POST("/verify-otp", authController.VerifyOtp)
	}
}
