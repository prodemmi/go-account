package container

import (
	"account/amqp"
	"account/http/controllers"
	"account/http/middlewares"
	"account/http/repositories"
	"account/http/services"
	"account/internal"

	"github.com/gin-gonic/gin"
)

type Container struct {
	AppRepository repositories.IAppRepository
	AppService    services.IAppService
	AppController controllers.IAppController

	UserRepository repositories.IUserRepository
	UserService    services.IUserService
	UserController controllers.IUserController

	OTPRepository repositories.IOTPRepository
	OTPService    services.IOTPService

	AuthRepository repositories.IAuthRepository
	AuthService    services.IAuthService
	AuthController controllers.IAuthController

	AppMiddleware  gin.HandlerFunc
	JWTMiddleware  gin.HandlerFunc
	CORSMiddleware gin.HandlerFunc

	AMQPService internal.IAMQPService
	AppAMQP     amqp.IAppAMQP
}

func NewContainer() *Container {
	appRepo := repositories.NewAppRepository()
	appServ := services.NewAppService(appRepo)
	appCont := controllers.NewAppController(appServ)

	userRepo := repositories.NewUserRepository()
	userServ := services.NewUserService(userRepo)
	userCont := controllers.NewUserController(userServ)

	otpRepo := repositories.NewOTPRepository()
	otpServ := services.NewOTPService(otpRepo)

	authRepo := repositories.NewAuthRepository()
	authServ := services.NewAuthService(authRepo, appServ, userServ)
	authCont := controllers.NewAuthController(authServ, userServ, otpServ)

	appMiddleware := middlewares.NewAppMiddleware(appServ)
	jwtMiddleware := middlewares.NewJWTMiddleware(userServ)
	corsMiddleware := middlewares.NewCORSMiddleware()

	amqpService := internal.NewAMQPService()
	appAMQP := amqp.NewAppAMQP(amqpService, authServ)

	return &Container{
		AppRepository: appRepo,
		AppService:    appServ,
		AppController: appCont,

		UserRepository: userRepo,
		UserService:    userServ,
		UserController: userCont,

		OTPRepository: otpRepo,
		OTPService:    otpServ,

		AuthRepository: authRepo,
		AuthService:    authServ,
		AuthController: authCont,

		AppMiddleware:  appMiddleware,
		JWTMiddleware:  jwtMiddleware,
		CORSMiddleware: corsMiddleware,

		AMQPService: amqpService,
		AppAMQP:     appAMQP,
	}
}
