package handlers

import (
	"fmt"
	"net/http"
	"time"

	"aman/internal/domains/otp"
	"aman/internal/domains/session"
	"aman/internal/ports/outbound"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler is the primary (inbound) HTTP adapter for CAS authentication.
type AuthHandler struct {
	appRepo       outbound.IAppRepository
	sessionRepo   outbound.ISessionRepository
	otpRepo       outbound.IOTPRepository
	authHash      outbound.IAuthHash
	jwt           outbound.IJWTService
	otpGen        outbound.IOTPGenerator
	notification  outbound.INotificationService
}

func NewAuthHandler(
	appRepo outbound.IAppRepository,
	sessionRepo outbound.ISessionRepository,
	otpRepo outbound.IOTPRepository,
	authHash outbound.IAuthHash,
	jwt outbound.IJWTService,
	otpGen outbound.IOTPGenerator,
	notification outbound.INotificationService,
) *AuthHandler {
	return &AuthHandler{
		appRepo:        appRepo,
		sessionRepo:    sessionRepo,
		otpRepo:        otpRepo,
		authHash:       authHash,
		jwt:            jwt,
		otpGen:         otpGen,
		notification:   notification,
	}
}

// SignIn godoc
// @Summary SignIn an app (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body signInBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/auth/sign-in [post]
func (h AuthHandler) SignIn(c *gin.Context) {
	var payload signInBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	appIDVal, _ := c.Get("AppID")
	appUUID, err := uuid.Parse(fmt.Sprintf("%v", appIDVal))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": "error",
			"error":  "invalid app id",
		})
		return
	}

	a, err := h.appRepo.FindByID(appUUID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "app not found.",
		})
		return
	}

	otpHash := h.authHash.CreateAuthHash(*a, c.Request.UserAgent())
	otpCode, err := h.otpGen.GenerateOTP()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	expireAt := time.Now().Add(2 * time.Minute)
	if _, err := h.otpRepo.Create(otp.Otp{
		AppID:    appUUID,
		Type:     a.OtpType,
		Code:     otpCode,
		Hash:     &otpHash,
		ExpireAt: &expireAt,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Send the OTP via SMS (notification adapter with error handling).
	if payload.Mobile != "" {
		if err := h.notification.SendSMS(payload.Mobile, "Your OTP code is "+otpCode); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"app_id":   a.ID,
			"otp_hash": otpHash,
		},
	})
}

// SignUp godoc
// @Summary SignUp a user (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body signUpBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/auth/sign-up [post]
func (h AuthHandler) SignUp(c *gin.Context) {
	var payload signUpBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	appIDVal, _ := c.Get("AppID")
	appUUID, err := uuid.Parse(fmt.Sprintf("%v", appIDVal))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": "error",
			"error":  "invalid app id",
		})
		return
	}

	a, err := h.appRepo.FindByID(appUUID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "app not found.",
		})
		return
	}

	otpHash := h.authHash.CreateAuthHash(*a, c.Request.UserAgent())
	otpCode, err := h.otpGen.GenerateOTP()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	expireAt := time.Now().Add(2 * time.Minute)
	if _, err := h.otpRepo.Create(otp.Otp{
		AppID:    appUUID,
		Type:     a.OtpType,
		Code:     otpCode,
		Hash:     &otpHash,
		ExpireAt: &expireAt,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if payload.Mobile != "" {
		if err := h.notification.SendSMS(payload.Mobile, "Your OTP code is "+otpCode); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"app_id":   a.ID,
			"otp_hash": otpHash,
		},
	})
}

// VerifyOtp godoc
// @Summary Verify otp and create session (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param payload body verifyOTPBody true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/auth/verify-otp [post]
func (h AuthHandler) VerifyOtp(c *gin.Context) {
	var payload verifyOTPBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	appIDVal, _ := c.Get("AppID")
	appUUID, err := uuid.Parse(fmt.Sprintf("%v", appIDVal))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "invalid app id",
		})
		return
	}

	a, err := h.appRepo.FindByID(appUUID.String())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "app not found.",
		})
		return
	}

	if !h.authHash.VerifyAuthHash(payload.OtpHash, *a, c.Request.UserAgent()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "cannot verify otp hash",
		})
		return
	}

	verified, err := h.otpRepo.CheckOTP(appUUID.String(), payload.OtpCode, payload.OtpType, payload.OtpHash)
	if err != nil || !verified {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "cannot verify otp",
		})
		return
	}

	token, err := h.jwt.CreateToken(appUUID.String(), a.SessionTimeout)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	expiresAt := time.Now().Add(a.SessionTimeout)
	if _, err := h.sessionRepo.Create(session.Session{
		AppID:     appUUID,
		Token:     token,
		ExpiresAt: expiresAt,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	appUrl, _ := c.Get("AppAuthUrl")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"api_token":   token,
			"redirect_to": appUrl,
		},
	})
}

// Me godoc
// @Summary Get authenticated app (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param API-TOKEN header string true "null"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/auth/me [get]
func (h AuthHandler) Me(c *gin.Context) {
	appEntity, exist := c.Get("App")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "app not authorized.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   appEntity,
	})
}

// OAuth2 is a placeholder for the OAuth2 flow.
func (h AuthHandler) OAuth2(c *gin.Context) {}

// OAuth2Fallback is a placeholder for the OAuth2 fallback flow.
func (h AuthHandler) OAuth2Fallback(c *gin.Context) {}

// Logout godoc
// @Summary Logout app (CAS)
// @Tags Auth
// @Accept json
// @Produce json
// @Param APP-ID header string true "null"
// @Param API-TOKEN header string true "null"
// @Success 308 {string} string "Location"
// @Router /v1/auth/logout [delete]
func (h AuthHandler) Logout(c *gin.Context) {
	appUrl, _ := c.Get("AppUrl")
	c.JSON(http.StatusPermanentRedirect, appUrl)
}
