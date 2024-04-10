package middlewares

import (
	"account/services"
	"account/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppHeader struct {
	AppID string `header:"APP-ID" binding:"required"`
}

func App(appService services.AppService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &AppHeader{}
		if err := c.BindHeader(header); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "app id is not set.",
			})
			c.Abort()
			return
		}

		if !utils.IsUUID(header.AppID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "app id is not valid uuid.",
			})
			c.Abort()
			return
		}

		app, err := appService.Exists(header.AppID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		if app {
			c.Next()
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "app not found,",
			})
		}
		c.Abort()
		return
	}
}
