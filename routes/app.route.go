package routes

import (
	"account/controllers"
	"account/repositories"
	"account/services"
	"github.com/gin-gonic/gin"
)

func RegisterAppRoutes(app *gin.Engine) {
	appRepository := repositories.NewAppRepository()
	appService := services.NewAppService(appRepository)
	appController := controllers.NewAppController(appService)

	v1 := app.Group("/v1/apps/")
	{
		v1.POST("/", appController.Store)
		v1.GET("/", appController.Index)
		v1.GET("/:app", appController.Single)
		v1.DELETE("/:id", appController.Delete)
	}
}
