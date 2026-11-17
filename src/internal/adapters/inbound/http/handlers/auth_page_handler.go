package handlers

import (
	"net/http"

	"aman/internal/adapters/inbound/http/pages"
	"aman/internal/adapters/inbound/http/pages/auth"
	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
)

// AuthPageHandler serves branded auth pages (CAS responsibility).
// CAS is the central authentication service: apps redirect users here for
// auth, and CAS renders app-branded login/signup/verify-OTP pages.
type AuthPageHandler struct {
	appRepo outbound.IAppRepository
}

func NewAuthPageHandler(appRepo outbound.IAppRepository) *AuthPageHandler {
	return &AuthPageHandler{appRepo: appRepo}
}

// Login renders the branded login page for an app.
// If the app has a custom AuthPageURL, redirect there instead.
func (h AuthPageHandler) Login(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "app_id is required"})
		return
	}

	a, err := h.appRepo.FindByID(appID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	// If app has a custom auth page URL, redirect there (app hosts its own auth)
	if a.AuthPageEnabled && a.AuthPageURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, a.AuthPageURL+"?app_id="+appID)
		return
	}

	// CAS serves the branded auth page
	primaryColor := a.AuthPagePrimaryColor
	if primaryColor == "" {
		primaryColor = "#3b82f6"
	}
	template := a.AuthPageTemplate
	if template == "" {
		template = "default"
	}

	component := auth.Login(pages.AuthData{
		App:             *a,
		AppID:           appID,
		PageTitle:       "Sign in",
		PrimaryColor:    primaryColor,
		LogoURL:         a.AuthPageLogo,
		CsrfToken:       generateCsrfToken(),
	})

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// Signup renders the branded signup page for an app.
func (h AuthPageHandler) Signup(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "app_id is required"})
		return
	}

	a, err := h.appRepo.FindByID(appID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	if a.AuthPageEnabled && a.AuthPageURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, a.AuthPageURL+"/signup?app_id="+appID)
		return
	}

	primaryColor := a.AuthPagePrimaryColor
	if primaryColor == "" {
		primaryColor = "#3b82f6"
	}

	component := auth.Signup(pages.AuthData{
		App:             *a,
		AppID:           appID,
		PageTitle:       "Create account",
		PrimaryColor:    primaryColor,
		LogoURL:         a.AuthPageLogo,
		CsrfToken:       generateCsrfToken(),
	})

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// VerifyOTP renders the branded verify-OTP page for an app.
func (h AuthPageHandler) VerifyOTP(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "app_id is required"})
		return
	}

	a, err := h.appRepo.FindByID(appID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	if a.AuthPageEnabled && a.AuthPageURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, a.AuthPageURL+"/verify-otp?app_id="+appID)
		return
	}

	primaryColor := a.AuthPagePrimaryColor
	if primaryColor == "" {
		primaryColor = "#3b82f6"
	}

	otpHash := c.Query("otp_hash")

	component := auth.VerifyOTP(pages.AuthData{
		App:        *a,
		AppID:      appID,
		PageTitle:  "Verify identity",
		PrimaryColor: primaryColor,
		LogoURL:    a.AuthPageLogo,
		CsrfToken:  generateCsrfToken(),
		OtpHash:    otpHash,
	})

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func generateCsrfToken() string {
	// Simple CSRF token; in production use crypto/rand
	return "csrf-placeholder"
}
