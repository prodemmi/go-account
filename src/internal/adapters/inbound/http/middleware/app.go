package middleware

import (
	"net/http"

	"aman/internal/adapters/inbound/http/handlers"
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
func NewAppMiddleware(appRepo outbound.AppRepository) AppMiddleware {
	return AppMiddleware(func(c *gin.Context) {
		header := &appHeader{}
		if err := c.BindHeader(header); err != nil {
			redirectNotFound(c)
			return
		}

		if _, err := uuid.Parse(header.AppID); err != nil {
			redirectNotFound(c)
			return
		}

		a, err := appRepo.FindByID(header.AppID)
		if err != nil {
			redirectNotFound(c)
			return
		}

		normalized := a.Normalize()
		c.Set("App", &normalized)
		c.Set("AppID", normalized.ID.String())
		c.Set("AppUrl", normalized.Url)
		c.Set("AppCallbackUrl", normalized.CallbackURL)
		c.Set("AppOtpType", normalized.OTPChannel)
		c.Next()
	})
}

func redirectNotFound(c *gin.Context) {
	c.Redirect(http.StatusFound, handlers.NotFoundPath)
	c.Abort()
}
