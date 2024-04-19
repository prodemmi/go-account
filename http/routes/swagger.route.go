package routes

import (
	"account/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerRoutes(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
