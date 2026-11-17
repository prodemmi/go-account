package admin

import (
	"net/http"
	"strconv"
	"time"

	"aman/internal/adapters/inbound/http/pages/admin"
	"aman/internal/domains/app"
	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
)

// AdminHandler manages apps and sessions through a separate admin panel.
// These routes are separate from the public API routes (/v1/*).
type AdminHandler struct {
	appRepo     outbound.IAppRepository
	sessionRepo outbound.ISessionRepository
}

func NewAdminHandler(
	appRepo outbound.IAppRepository,
	sessionRepo outbound.ISessionRepository,
) *AdminHandler {
	return &AdminHandler{
		appRepo:     appRepo,
		sessionRepo: sessionRepo,
	}
}

// AppsList renders the admin apps list.
func (h AdminHandler) AppsList(c *gin.Context) {
	apps, err := h.appRepo.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows := make([]admin.AppRow, len(apps))
	for i, a := range apps {
		rows[i] = admin.AppRow{
			ID:              a.ID.String(),
			Name:            a.Name,
			Url:             a.Url,
			OtpType:         a.OtpType,
			SessionTimeout:  int(a.SessionTimeout.Seconds()),
			AuthPageEnabled: a.AuthPageEnabled,
		}
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := admin.AppsList(rows).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// NewApp renders the new app form.
func (h AdminHandler) NewApp(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := admin.AppForm(admin.AppFormData{
		AuthPagePrimaryColor: "#3b82f6",
		AuthPageTemplate:     "default",
	}, false).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// CreateApp creates a new app from the admin form.
func (h AdminHandler) CreateApp(c *gin.Context) {
	var payload appFormPayload
	if err := c.ShouldBind(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionTimeout, _ := strconv.Atoi(payload.SessionTimeout)
	authPageEnabled := payload.AuthPageEnabled == "true"

	a, err := h.appRepo.Create(app.App{
		Name:                 payload.Name,
		Url:                  payload.Url,
		AuthUrl:              payload.AuthUrl,
		OtpType:              payload.OtpType,
		SessionTimeout:       time.Duration(sessionTimeout) * time.Second,
		AuthPageURL:          payload.AuthPageURL,
		AuthPageTemplate:     payload.AuthPageTemplate,
		AuthPageLogo:         payload.AuthPageLogo,
		AuthPagePrimaryColor: payload.AuthPagePrimaryColor,
		AuthPageEnabled:      authPageEnabled,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/apps/"+a.ID.String())
}

// EditApp renders the edit app form.
func (h AdminHandler) EditApp(c *gin.Context) {
	id := c.Param("id")

	a, err := h.appRepo.FindByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := admin.AppForm(admin.AppFormData{
		ID:                   a.ID.String(),
		Name:                 a.Name,
		Url:                  a.Url,
		AuthUrl:              a.AuthUrl,
		OtpType:              a.OtpType,
		SessionTimeout:       int(a.SessionTimeout.Seconds()),
		AuthPageURL:          a.AuthPageURL,
		AuthPageTemplate:     a.AuthPageTemplate,
		AuthPageLogo:         a.AuthPageLogo,
		AuthPagePrimaryColor: a.AuthPagePrimaryColor,
		AuthPageEnabled:      a.AuthPageEnabled,
	}, true).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// UpdateApp updates an existing app.
func (h AdminHandler) UpdateApp(c *gin.Context) {
	id := c.Param("id")

	var payload appFormPayload
	if err := c.ShouldBind(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionTimeout, _ := strconv.Atoi(payload.SessionTimeout)
	authPageEnabled := payload.AuthPageEnabled == "true"

	a, err := h.appRepo.FindByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	a.Name = payload.Name
	a.Url = payload.Url
	a.AuthUrl = payload.AuthUrl
	a.OtpType = payload.OtpType
	a.SessionTimeout = time.Duration(sessionTimeout) * time.Second
	a.AuthPageURL = payload.AuthPageURL
	a.AuthPageTemplate = payload.AuthPageTemplate
	a.AuthPageLogo = payload.AuthPageLogo
	a.AuthPagePrimaryColor = payload.AuthPagePrimaryColor
	a.AuthPageEnabled = authPageEnabled

	if _, err := h.appRepo.Update(*a); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/apps")
}

// DeleteApp deletes an app.
func (h AdminHandler) DeleteApp(c *gin.Context) {
	id := c.Param("id")

	if err := h.appRepo.Delete(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/apps")
}

// SessionsList renders the admin sessions list.
func (h AdminHandler) SessionsList(c *gin.Context) {
	sessions, err := h.sessionRepo.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows := make([]admin.SessionRow, len(sessions))
	for i, s := range sessions {
		rows[i] = admin.SessionRow{
			AppID:     s.AppID.String(),
			Token:     s.Token,
			ExpiresAt: s.ExpiresAt.String(),
		}
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := admin.SessionsList(rows).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// DeleteSession revokes a session.
func (h AdminHandler) DeleteSession(c *gin.Context) {
	token := c.Param("token")

	if err := h.sessionRepo.Delete(token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/sessions")
}

type appFormPayload struct {
	Name                 string `form:"name" binding:"required"`
	Url                  string `form:"url" binding:"required"`
	AuthUrl              string `form:"auth_url" binding:"required"`
	OtpType              string `form:"otp_type"`
	SessionTimeout       string `form:"session_timeout"`
	AuthPageURL          string `form:"auth_page_url"`
	AuthPageTemplate     string `form:"auth_page_template"`
	AuthPageLogo         string `form:"auth_page_logo"`
	AuthPagePrimaryColor string `form:"auth_page_primary_color"`
	AuthPageEnabled      string `form:"auth_page_enabled"`
}
