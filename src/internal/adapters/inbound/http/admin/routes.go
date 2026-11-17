package admin

import (
	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes registers the admin panel routes (separate from /v1/*).
// Admin routes: GET/POST/PUT/DELETE /admin/apps, GET/POST /admin/sessions.
func RegisterAdminRoutes(
	router *gin.Engine,
	handler *AdminHandler,
) {
	adminGroup := router.Group("/admin")
	{
		// Apps
		adminGroup.GET("/apps", handler.AppsList)
		adminGroup.GET("/apps/new", handler.NewApp)
		adminGroup.POST("/apps", handler.CreateApp)
		adminGroup.GET("/apps/:id", handler.EditApp)
		adminGroup.POST("/apps/:id/update", handler.UpdateApp)
		adminGroup.POST("/apps/:id/delete", handler.DeleteApp)

		// Sessions
		adminGroup.GET("/sessions", handler.SessionsList)
		adminGroup.POST("/sessions/:token/delete", handler.DeleteSession)
	}
}
