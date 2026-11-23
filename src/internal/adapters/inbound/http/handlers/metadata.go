package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// requestMeta is the request fingerprint stored on a session. Country, city
// and region come from common reverse-proxy/geo headers (Cloudflare, nginx
// GeoIP modules) and are empty when none are set.
type requestMeta struct {
	UserAgent string
	IP        string
	Country   string
	City      string
	Region    string
	Device    string
}

func metaFrom(c *gin.Context) requestMeta {
	agent := c.Request.UserAgent()
	return requestMeta{
		UserAgent: agent,
		IP:        c.ClientIP(),
		Country:   firstHeader(c, "CF-IPCountry", "X-Country-Code", "X-Geo-Country"),
		City:      firstHeader(c, "CF-IPCity", "X-City", "X-Geo-City"),
		Region:    firstHeader(c, "CF-Region", "X-Region", "X-Geo-Region"),
		Device:    deviceFromAgent(agent),
	}
}

func firstHeader(c *gin.Context, names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(c.GetHeader(name)); value != "" {
			return value
		}
	}
	return ""
}

func deviceFromAgent(agent string) string {
	lower := strings.ToLower(agent)
	switch {
	case strings.Contains(lower, "ipad"), strings.Contains(lower, "tablet"):
		return "tablet"
	case strings.Contains(lower, "mobile"), strings.Contains(lower, "iphone"), strings.Contains(lower, "android"):
		return "mobile"
	case strings.Contains(lower, "bot"), strings.Contains(lower, "crawler"), strings.Contains(lower, "spider"):
		return "bot"
	case agent == "":
		return ""
	default:
		return "desktop"
	}
}
