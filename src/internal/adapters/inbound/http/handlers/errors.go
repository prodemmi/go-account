package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// NotFoundPath is served by nginx (see nginx.conf): the app never renders the
// error page itself, it redirects there instead of returning JSON.
const NotFoundPath = "/not-found"

func queryEscape(s string) string {
	return url.QueryEscape(s)
}

// redirectNotFound sends the browser to the nginx-defined 404 page.
func redirectNotFound(c *gin.Context) {
	c.Redirect(http.StatusFound, NotFoundPath)
	c.Abort()
}

// redirectWithError sends the user back to an auth page with a readable error.
func redirectWithError(c *gin.Context, target, message string) {
	sep := "?"
	if strings.Contains(target, "?") {
		sep = "&"
	}
	c.Redirect(http.StatusFound, target+sep+"error="+url.QueryEscape(message))
	c.Abort()
}
