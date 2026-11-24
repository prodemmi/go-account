package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aman/internal/domains/app"
	"aman/internal/ports/inbound"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type checkAuthStub struct {
	inbound.AuthUseCase
	command inbound.CheckCommand
	result  *inbound.CheckResult
}

func (s *checkAuthStub) Check(_ context.Context, command inbound.CheckCommand) (*inbound.CheckResult, error) {
	s.command = command
	if s.result != nil {
		return s.result, nil
	}
	return &inbound.CheckResult{
		Authenticated: true,
		OTT:           "one-time-ticket",
		RedirectURL:   "https://app.example/callback?ott=one-time-ticket",
	}, nil
}

func TestCheckLoginReturnsJSONAndAcceptsAPIToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	appID := "ef84cb80-4e0a-4995-b7ea-9537b4833b45"
	auth := &checkAuthStub{}
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("App", &app.App{ID: uuid.MustParse(appID)})
		c.Next()
	})
	router.GET("/v1/auth/check-login", NewAuthHandler(auth).CheckLogin)
	request := httptest.NewRequest(http.MethodGet, "/v1/auth/check-login", nil)
	request.Header.Set("API-TOKEN", "session-token")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if auth.command.SessionToken != "session-token" {
		t.Fatalf("session token = %q, want API-TOKEN header", auth.command.SessionToken)
	}
	var body struct {
		Status string `json:"status"`
		Data   struct {
			Authenticated bool   `json:"authenticated"`
			OTT           string `json:"ott"`
			RedirectURL   string `json:"redirect_url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "success" || !body.Data.Authenticated || body.Data.OTT != "one-time-ticket" {
		t.Fatalf("unexpected check-login response: %+v", body)
	}
}

func TestCustomModeCheckReturnsUnauthenticatedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("App", &app.App{ID: uuid.New(), CustomAuthEnabled: true})
		c.Next()
	})
	router.GET("/v1/auth/check", NewAuthHandler(&checkAuthStub{result: &inbound.CheckResult{}}).Check)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/auth/check", nil))

	var body struct {
		Data struct {
			Authenticated bool `json:"authenticated"`
		} `json:"data"`
	}
	if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &body) != nil || body.Data.Authenticated {
		t.Fatalf("custom check should return unauthenticated JSON: status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestAPISendOTPRejectsEnabledCaptchaWithoutKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("App", &app.App{
			ID:                      uuid.New(),
			AuthPageCaptchaEnabled:  true,
			AuthPageCaptchaProvider: app.CaptchaProviderGoogle,
		})
		c.Next()
	})
	router.POST("/v1/auth/send-otp", NewAuthHandler(&checkAuthStub{}).SendOTP)
	request := httptest.NewRequest(http.MethodPost, "/v1/auth/send-otp", strings.NewReader(`{"identifier":"user@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusFound || response.Header().Get("Location") != "/not-found" {
		t.Fatalf("status/location = %d/%q, want redirect to /not-found", response.Code, response.Header().Get("Location"))
	}
}
