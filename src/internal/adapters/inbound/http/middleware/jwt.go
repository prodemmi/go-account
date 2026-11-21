package middleware

import (
	"fmt"
	"net/http"
	"time"

	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a named gin.HandlerFunc type so fx can distinguish it
// from other handler funcs (e.g. AppMiddleware).
type JWTMiddleware gin.HandlerFunc

// NewJWTMiddleware validates the API-TOKEN header using the JWT outbound port
// and stores the app ID on the gin context.
func NewJWTMiddleware(jwt outbound.IJWTService, sessionRepo outbound.ISessionRepository) JWTMiddleware {
	return JWTMiddleware(func(c *gin.Context) {
		token := c.GetHeader("API-TOKEN")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "api token not provided",
			})
			return
		}

		appID, err := jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "invalid api token",
			})
			return
		}

		if headerAppID, _ := c.Get("AppID"); fmt.Sprintf("%v", headerAppID) != appID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "token app mismatch",
			})
			return
		}

		s, err := sessionRepo.FindByToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "session not found",
			})
			return
		}

		if s.AppID.String() != appID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "session app mismatch",
			})
			return
		}

		if s.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "session expired",
			})
			return
		}

		c.Set("User", appID)
		c.Set("UserID", appID)
		c.Set("AppID", appID)
		c.Next()
	})
}
