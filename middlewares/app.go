package middlewares

import (
	"github.com/gin-gonic/gin"
	"account/services"
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
				"error": err.Error(),
			})
			return
		}

		if app, err := appService.Exists(header.AppID); err != nil {
			if app {
				c.Next()
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "app not found,",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}
}
