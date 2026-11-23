package middleware

import (
	"time"

	"aman/internal/domains/app"
	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a named gin.HandlerFunc type so fx can distinguish it
// from other handler funcs (e.g. AppMiddleware).
type JWTMiddleware gin.HandlerFunc

// NewJWTMiddleware validates the session token (API-TOKEN header or the CAS
// cookie), loads the session and stores it on the gin context.
func NewJWTMiddleware(jwt outbound.JWTService, sessionRepo outbound.SessionRepository) JWTMiddleware {
	return JWTMiddleware(func(c *gin.Context) {
		token := c.GetHeader("API-TOKEN")
		if token == "" {
			token = cookieToken(c)
		}
		if token == "" {
			redirectNotFound(c)
			return
		}

		appID, identifier, err := jwt.ParseToken(token)
		if err != nil {
			redirectNotFound(c)
			return
		}

		// The token must belong to the app named in APP-ID; otherwise a
		// session from one app could be used against another app's routes.
		if requestAppID := c.GetString("AppID"); requestAppID != "" && requestAppID != appID {
			redirectNotFound(c)
			return
		}

		s, err := sessionRepo.FindByToken(token)
		if err != nil || s.AppID.String() != appID || !s.Active(time.Now()) || !s.MatchesAgent(c.Request.UserAgent()) {
			redirectNotFound(c)
			return
		}

		c.Set("Session", s)
		c.Set("SessionID", s.ID.String())
		c.Set("User", identifier)
		c.Set("UserID", identifier)
		c.Set("Identifier", identifier)
		c.Set("AppID", appID)
		c.Next()
	})
}

// cookieToken reads the app-specific session cookie, then the defaults so a
// cookie rename does not immediately invalidate existing sessions.
func cookieToken(c *gin.Context) string {
	names := []string{app.DefaultSessionCookieName, "cas_session"}
	if value, exists := c.Get("App"); exists {
		if a, ok := value.(*app.App); ok && a != nil {
			names = append([]string{a.Normalize().SessionCookieName}, names...)
		}
	}
	seen := map[string]bool{}
	for _, name := range names {
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		if token, err := c.Cookie(name); err == nil && token != "" {
			return token
		}
	}
	return ""
}
