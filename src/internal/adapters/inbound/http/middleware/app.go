package middleware

import (
	"net/http"

	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AppMiddleware is a named gin.HandlerFunc type so fx can distinguish it
// from other handler funcs (e.g. JWTMiddleware).
type AppMiddleware gin.HandlerFunc

type appHeader struct {
	AppID string `header:"APP-ID" binding:"required"`
}

// NewAppMiddleware validates the APP-ID header, resolves the app and stores
// app-scoped values on the gin context for downstream handlers.
func NewAppMiddleware(appRepo outbound.IAppRepository) AppMiddleware {
	return AppMiddleware(func(c *gin.Context) {
		header := &appHeader{}
		if err := c.BindHeader(header); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "app id is not set.",
			})
			c.Abort()
			return
		}

		if _, err := uuid.Parse(header.AppID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "app id is not valid uuid.",
			})
			c.Abort()
			return
		}

		a, err := appRepo.FindByID(header.AppID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "app not found,",
			})
			c.Abort()
			return
		}

		c.Set("App", a)
		c.Set("AppID", a.ID)
		c.Set("AppUrl", a.Url)
		c.Set("AppAuthUrl", a.AuthUrl)
		c.Set("AppOtpType", a.OtpType)
		c.Next()
	})
}
