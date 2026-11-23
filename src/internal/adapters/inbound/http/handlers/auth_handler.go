package handlers

import (
	"net/http"

	"aman/internal/domains/app"
	"aman/internal/domains/session"
	"aman/internal/ports/inbound"

	"github.com/gin-gonic/gin"
)

// AuthHandler is the JSON API for CAS authentication. It never returns JSON
// on errors: invalid requests are redirected to the nginx 404 page.
type AuthHandler struct {
	auth inbound.AuthUseCase
}

func NewAuthHandler(auth inbound.AuthUseCase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// SendOTP creates and delivers an OTP code.
// @Summary Send OTP (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "app id"
// @Param payload body sendOTPBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/send-otp [post]
func (h AuthHandler) SendOTP(c *gin.Context) {
	a, appID, ok := requestApp(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	var payload sendOTPBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		redirectNotFound(c)
		return
	}

	meta := metaFrom(c)
	challenge, err := h.auth.SendOTP(c.Request.Context(), inbound.SendOTPCommand{
		AppID:      appID,
		Identifier: payload.Identifier,
		UserAgent:  meta.UserAgent,
		IP:         meta.IP,
		Country:    meta.Country,
		City:       meta.City,
		Region:     meta.Region,
		Device:     meta.Device,
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"app_id":     a.ID,
			"identifier": challenge.Identifier,
			"channel":    challenge.Channel,
			"otp_hash":   challenge.OTPHash,
			"length":     challenge.Length,
			"expires_at": challenge.ExpiresAt,
		},
	})
}

// VerifyOTP verifies a code, creates the pending CAS session and issues an OTT.
// @Summary Verify OTP (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "app id"
// @Param payload body verifyOTPBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/verify-otp [post]
func (h AuthHandler) VerifyOTP(c *gin.Context) {
	var payload verifyOTPBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		redirectNotFound(c)
		return
	}

	_, appID, ok := requestApp(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	meta := metaFrom(c)
	result, err := h.auth.VerifyOTP(c.Request.Context(), inbound.VerifyOTPCommand{
		AppID:      appID,
		Identifier: payload.Identifier,
		Code:       payload.OTPCode,
		OTPHash:    payload.OTPHash,
		UserAgent:  meta.UserAgent,
		IP:         meta.IP,
		Country:    meta.Country,
		City:       meta.City,
		Region:     meta.Region,
		Device:     meta.Device,
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	setSessionCookie(c, result.App, result.Session.Token, result.Session.ExpiresAt, result.Session.Infinite)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"session_token": result.Session.Token,
			"ott":           result.OTT,
			"redirect_url":  result.RedirectURL,
			"identifier":    result.Session.Identifier,
			"status":        result.Session.Status,
			"expires_at":    result.Session.ExpiresAt,
			"infinite":      result.Session.Infinite,
		},
	})
}

// VerifyOTT exchanges a one-time ticket for credentials.
// @Summary Verify OTT (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "app id"
// @Param payload body verifyOTTBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/verify-ott [post]
func (h AuthHandler) VerifyOTT(c *gin.Context) {
	var payload verifyOTTBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		redirectNotFound(c)
		return
	}

	a, appID, ok := requestApp(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	meta := metaFrom(c)
	credentials, err := h.auth.VerifyOTT(c.Request.Context(), inbound.VerifyOTTCommand{
		AppID:     appID,
		OTT:       payload.OTT,
		UserAgent: meta.UserAgent,
		IP:        meta.IP,
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	setSessionCookie(c, a, credentials.SessionToken, credentials.ExpiresAt, credentials.Infinite)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   credentials,
	})
}

// Check reports whether the browser already has a valid CAS session. When it
// does the user is redirected back to the app with a fresh OTT; otherwise the
// user goes to the auth page.
// @Summary Check CAS session (CAS)
// @Tags Auth
// @Param APP-ID header string true "app id"
// @Success 302 {string} string "Location"
// @Router /v1/auth/check [get]
func (h AuthHandler) Check(c *gin.Context) {
	a, appID, ok := requestApp(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	result, err := h.auth.Check(c.Request.Context(), inbound.CheckCommand{
		AppID:        appID,
		SessionToken: sessionCookie(c, a),
		UserAgent:    c.Request.UserAgent(),
		ReturnTo:     c.Query("redirect"),
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	if result.Authenticated {
		c.Redirect(http.StatusFound, result.RedirectURL)
		return
	}

	c.Redirect(http.StatusFound, loginPageURL(a, appID))
}

// Me returns the current CAS session.
// @Summary Current session (CAS)
// @Tags Auth
// @Param APP-ID header string true "app id"
// @Param API-TOKEN header string false "session token"
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/me [get]
func (h AuthHandler) Me(c *gin.Context) {
	sess, ok := currentSession(c)
	if !ok {
		redirectNotFound(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   sess,
	})
}

// Logout revokes the current CAS session.
// @Summary Logout (CAS)
// @Tags Auth
// @Param APP-ID header string true "app id"
// @Success 302 {string} string "Location"
// @Router /v1/auth/logout [delete]
func (h AuthHandler) Logout(c *gin.Context) {
	a, appID, ok := requestApp(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	token := sessionCookie(c, a)
	if token == "" {
		token = c.GetHeader("API-TOKEN")
	}

	if err := h.auth.Logout(c.Request.Context(), inbound.LogoutCommand{AppID: appID, SessionToken: token}); err != nil {
		redirectNotFound(c)
		return
	}

	clearSessionCookies(c, a)
	c.Redirect(http.StatusFound, a.Url)
}

func requestApp(c *gin.Context) (*app.App, string, bool) {
	value, exists := c.Get("App")
	if !exists {
		return nil, "", false
	}
	a, ok := value.(*app.App)
	if !ok || a == nil {
		return nil, "", false
	}
	return a, a.ID.String(), true
}

func currentSession(c *gin.Context) (*session.Session, bool) {
	value, exists := c.Get("Session")
	if !exists {
		return nil, false
	}
	sess, ok := value.(*session.Session)
	return sess, ok && sess != nil
}

func loginPageURL(a *app.App, appID string) string {
	target := "/auth/login?app_id=" + appID
	if a.CallbackURL != "" {
		target += "&redirect=" + queryEscape(a.CallbackURL)
	}
	return target
}
