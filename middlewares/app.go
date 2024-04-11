package middlewares

import (
	"account/repositories"
	"account/services"
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

		app, err := appService.CheckExists(header.AppID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		if app {
			c.Set("App", app)
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

func NewAppMiddleware() gin.HandlerFunc {
	appRepository := repositories.NewAppRepository()
	appService := services.NewAppService(appRepository)

	return App(appService)
}
