package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"aman/internal/domains/app"

	"github.com/gin-gonic/gin"
)

// SessionCookieName is the legacy default cookie name. Each app can override
// it in the admin panel; cas_cookie is the current default.
const SessionCookieName = app.DefaultSessionCookieName

func cookieName(a *app.App) string {
	if a == nil {
		return SessionCookieName
	}
	normalized := a.Normalize()
	if normalized.SessionCookieName == "" {
		return SessionCookieName
	}
	return normalized.SessionCookieName
}

func setSessionCookie(c *gin.Context, a *app.App, token string, expiresAt time.Time, infinite bool) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if infinite || maxAge <= 0 {
		maxAge = 10 * 365 * 24 * 60 * 60
	}

	normalized := app.App{}
	if a != nil {
		normalized = a.Normalize()
	}

	cookie := &http.Cookie{
		Name:     cookieName(a),
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Expires:  expiresAt,
		HttpOnly: normalized.SessionCookieHTTPOnly,
		Secure:   normalized.SessionCookieSecure || isSecureRequest(c),
		SameSite: sameSiteMode(normalized.SessionCookieSameSite),
	}

	for _, domain := range cookieDomains(c, a) {
		ck := *cookie
		ck.Domain = domain
		http.SetCookie(c.Writer, &ck)
	}
}

func clearSessionCookies(c *gin.Context, a *app.App) {
	normalized := app.App{}
	if a != nil {
		normalized = a.Normalize()
	}

	cookie := &http.Cookie{
		Name:     cookieName(a),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: normalized.SessionCookieHTTPOnly,
		Secure:   normalized.SessionCookieSecure || isSecureRequest(c),
		SameSite: sameSiteMode(normalized.SessionCookieSameSite),
	}
	for _, domain := range cookieDomains(c, a) {
		ck := *cookie
		ck.Domain = domain
		http.SetCookie(c.Writer, &ck)
	}
}

// sessionCookie reads the session token. It first tries the app-specific
// cookie name and then the historical default so sessions created before a
// rename keep working.
func sessionCookie(c *gin.Context, a *app.App) string {
	names := []string{cookieName(a), SessionCookieName, "cas_session"}
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

func sameSiteMode(value string) http.SameSite {
	switch value {
	case app.SameSiteStrict:
		return http.SameSiteStrictMode
	case app.SameSiteNone:
		return http.SameSiteNoneMode
	case app.SameSiteDefault:
		return http.SameSiteDefaultMode
	default:
		return http.SameSiteLaxMode
	}
}

// cookieDomains returns the domains the CAS cookie should be written to. CAS
// can only write cookies for its own host and for a shared parent domain; when
// the app lives on an unrelated domain the OTT redirect is the only bridge.
func cookieDomains(c *gin.Context, a *app.App) []string {
	casHost := hostOnly(c.Request.Host)
	if a == nil {
		return []string{""}
	}
	appHost := hostOnly(a.CallbackURL)
	if appHost == "" {
		appHost = hostOnly(a.Url)
	}
	if appHost == "" || appHost == casHost {
		return []string{""}
	}

	if parent := sharedParentDomain(casHost, appHost); parent != "" {
		return []string{parent}
	}

	// Best effort: the browser will ignore a domain it does not belong to.
	return []string{"", appHost}
}

func isSecureRequest(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}

func hostOnly(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if u, err := url.Parse(raw); err == nil && u.Hostname() != "" {
		return strings.ToLower(u.Hostname())
	}
	// raw is already a host:port.
	if i := strings.LastIndex(raw, ":"); i >= 0 {
		return strings.ToLower(raw[:i])
	}
	return strings.ToLower(raw)
}

func sharedParentDomain(a, b string) string {
	ap := strings.Split(a, ".")
	bp := strings.Split(b, ".")
	i, j := len(ap)-1, len(bp)-1
	shared := make([]string, 0, 3)
	for i >= 0 && j >= 0 && ap[i] == bp[j] {
		shared = append([]string{ap[i]}, shared...)
		i--
		j--
	}
	if len(shared) < 2 {
		return ""
	}
	return strings.Join(shared, ".")
}
