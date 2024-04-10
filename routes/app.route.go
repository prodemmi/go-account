package routes

import (
	"github.com/gin-gonic/gin"
	"account/controllers"
	"account/services"
	"gorm.io/gorm"
)

func RegisterAppRoutes(gin *gin.Engine, db *gorm.DB) {
	appController := controllers.AppController{
		App: gin,
		AppService: services.AppService{
			DB:  db,
			App: gin,
		},
	}

	v1 := gin.Group("/v1/apps/")
	{
		v1.POST("/", appController.Store)
		v1.GET("/", appController.Index)
		v1.GET("/:app", appController.Single)
	}
}
