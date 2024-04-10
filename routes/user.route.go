package routes

import (
	"github.com/gin-gonic/gin"
	"account/controllers"
	"account/middlewares"
	"account/services"
	"gorm.io/gorm"
)

func RegisterUserRoutes(gin *gin.Engine, db *gorm.DB) {
	appService := *services.CreateNewAppService(gin, db)
	userController := controllers.UserController{
		App: gin,
		UserService: services.UserService{
			App:          gin,
			DB:           db,
			TokenService: *services.CreateNewTokenService(gin, db),
		},
		AppService: appService,
	}

	v1 := gin.Group("/v1/users/").Use(middlewares.App(appService))
	{
		v1.GET("/", userController.Index)
		v1.GET("/:id", userController.Single)

		v1.POST("/", userController.Store)
	}
}
