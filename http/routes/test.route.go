package routes

import (
	"account/http/container"
	"account/internal"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
)

func RegisterTestRoutes(app *gin.Engine, container *container.Container) {
	v1 := app.Group("/v1/test/").Use(container.AppMiddleware, container.JWTMiddleware)
	{
		v1.GET("/sms", func(context *gin.Context) {
			a, _ := context.Get("AppID")
			t, _ := context.Get("Token")

			appID := utils.ToString(a)
			token := utils.ToString(t)
			container.AMQPService.Publish(appID, internal.Message{
				AppID:    &appID,
				ApiToken: &token,
				Event:    "SEND_SMS",
				Body:     "09130243519",
			})
		})
		v1.GET("/email", func(context *gin.Context) {
			a, _ := context.Get("AppID")
			t, _ := context.Get("Token")

			appID := utils.ToString(a)
			token := utils.ToString(t)
			container.AMQPService.Publish(appID, internal.Message{
				AppID:    &appID,
				ApiToken: &token,
				Event:    "SEND_EMAIL",
				Body:     "prodemmi@gmail.com",
			})
		})
	}
}
