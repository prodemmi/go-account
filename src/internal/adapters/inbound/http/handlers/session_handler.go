package handlers

import (
	"net/http"
	"time"

	"aman/internal/domains/session"
	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SessionHandler is the primary (inbound) HTTP adapter for session management (CAS).
type SessionHandler struct {
	sessionRepo outbound.ISessionRepository
	appRepo     outbound.IAppRepository
	jwt         outbound.IJWTService
}

func NewSessionHandler(
	sessionRepo outbound.ISessionRepository,
	appRepo outbound.IAppRepository,
	jwt outbound.IJWTService,
) *SessionHandler {
	return &SessionHandler{
		sessionRepo: sessionRepo,
		appRepo:     appRepo,
		jwt:         jwt,
	}
}

// Create godoc
// @Summary Create a new CAS session
// @Tags Session
// @Accept json
// @Produce json
// @Param payload body createSessionBody true "payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/sessions [post]
func (h SessionHandler) Create(c *gin.Context) {
	var payload createSessionBody
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	appID, err := uuid.Parse(payload.AppID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "invalid app id",
		})
		return
	}

	a, err := h.appRepo.FindByID(payload.AppID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "app not found.",
		})
		return
	}

	token, err := h.jwt.CreateToken(payload.AppID, a.SessionTimeout)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	s, err := h.sessionRepo.Create(session.Session{
		AppID:     appID,
		Token:     token,
		ExpiresAt: time.Now().Add(a.SessionTimeout),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   s,
	})
}

// Validate godoc
// @Summary Validate a CAS session token
// @Tags Session
// @Accept json
// @Produce json
// @Param token query string true "session token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/sessions/validate [get]
func (h SessionHandler) Validate(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "token is required",
		})
		return
	}

	appID, err := h.jwt.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "invalid token",
		})
		return
	}

	s, err := h.sessionRepo.FindByToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "session not found",
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"app_id":   appID,
			"session":  s,
			"valid":    true,
		},
	})
}

// Delete godoc
// @Summary Delete a CAS session (logout)
// @Tags Session
// @Accept json
// @Produce json
// @Param token query string true "session token"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/sessions [delete]
func (h SessionHandler) Delete(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "token is required",
		})
		return
	}

	if err := h.sessionRepo.Delete(token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   true,
	})
}
