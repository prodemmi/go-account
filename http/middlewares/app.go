package middlewares

import (
	"account/http/services"
	"account/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppHeader struct {
	AppID string `header:"APP-ID" binding:"required"`
}

func App(appService services.IAppService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &AppHeader{}
		if err := c.BindHeader(header); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "server id is not set.",
			})
			c.Abort()
			return
		}

		if !utils.IsUUID(header.AppID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "server id is not valid uuid.",
			})
			c.Abort()
			return
		}

		app, err := appService.CheckExists(header.AppID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		if app {
			appEntity, _ := appService.GetApp(header.AppID)
			c.Set("App", appEntity)
			c.Set("AppID", appEntity.ID)
			c.Set("AppUrl", appEntity.Url)
			c.Set("AppOtpType", appEntity.OtpType)
			c.Next()
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "server not found,",
			})
		}
		c.Abort()
		return
	}
}

func NewAppMiddleware(appService services.IAppService) gin.HandlerFunc {
	return App(appService)
}
