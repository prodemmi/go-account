package admin

import (
	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes registers the server-rendered admin panel.
func RegisterAdminRoutes(router *gin.Engine, handler *AdminHandler) {
	group := router.Group("/admin")
	{
		group.GET("", handler.Dashboard)
		group.GET("/", handler.Dashboard)
		group.GET("/dashboard", handler.Dashboard)

		group.GET("/apps", handler.AppsList)
		group.GET("/apps/new", handler.NewApp)
		group.POST("/apps", handler.CreateApp)
		group.GET("/apps/:id", handler.EditApp)
		group.POST("/apps/:id/update", handler.UpdateApp)
		group.POST("/apps/:id/delete", handler.DeleteApp)

		group.GET("/sessions", handler.SessionsList)
		group.POST("/sessions/:id/delete", handler.DeleteSession)

		group.GET("/guide", handler.Guide)
		group.GET("/developers", handler.Developers)
	}
}
