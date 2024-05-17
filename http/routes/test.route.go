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

			sender, _ := container.UserService.GetUserByToken(token)

			container.AMQPService.Publish(appID, internal.Message{
				AppID:    &appID,
				ApiToken: &token,
				Event:    "SEND_SMS",
				Body: internal.SmsBody{
					App:     appID,
					Sender:  sender.ID.String(),
					Message: "Hello emad",
					TO:      "09130243519",
				},
			})
		})
		v1.GET("/email", func(context *gin.Context) {
			a, _ := context.Get("AppID")
			t, _ := context.Get("Token")

			appID := utils.ToString(a)
			token := utils.ToString(t)

			sender, _ := container.UserService.GetUserByToken(token)

			container.AMQPService.Publish(appID, internal.Message{
				AppID:    &appID,
				ApiToken: &token,
				Event:    "SEND_EMAIL",
				Body: internal.EmailBody{
					App:     appID,
					Sender:  sender.ID.String(),
					Subject: "My email",
					CC:      "",
					BCC:     "",
					Message: "Hello emad",
					TO:      "prodemmi@gmail.com",
				},
			})
		})
	}
}
