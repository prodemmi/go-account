package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"aman/internal/domains/app"
	"aman/internal/ports/inbound"

	"github.com/gin-gonic/gin"
)

func TestIsSafeRedirectURL(t *testing.T) {
	for raw, want := range map[string]bool{
		"/login":                         true,
		"https://auth.example/login":     true,
		"//evil.example/login":           false,
		`\\evil.example\login`:           false,
		"/\\evil.example/login":          false,
		"javascript:alert(1)":            false,
		"https:/missing-host.example":    false,
		"https://user:pass@auth.example": false,
	} {
		if got := isSafeRedirectURL(raw); got != want {
			t.Errorf("isSafeRedirectURL(%q) = %t, want %t", raw, got, want)
		}
	}
}

type authPageAppStub struct {
	inbound.AppUseCase
	value *app.App
}

func (s authPageAppStub) Get(_ context.Context, _ string) (*app.App, error) {
	return s.value, nil
}

func TestRequireManagedAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		name       string
		app        *app.App
		url        string
		wantStatus int
	}{
		{name: "managed mode passes", app: &app.App{}, url: "/auth/login?app_id=app", wantStatus: http.StatusNoContent},
		{name: "custom mode is hidden", app: &app.App{CustomAuthEnabled: true}, url: "/auth/login?app_id=app", wantStatus: http.StatusNotFound},
		{name: "missing app is hidden", app: &app.App{}, url: "/auth/login", wantStatus: http.StatusNotFound},
	} {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			handler := AuthPageHandler{apps: authPageAppStub{value: test.app}}
			router.GET("/auth/login", handler.RequireManagedAuth(), func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.url, nil))
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
		})
	}
}
