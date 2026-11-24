package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"aman/internal/domains/app"
)

func TestGoogleRecaptchaProviderVerifiesServerSide(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected request: %s %s (%s)", r.Method, r.URL.Path, r.Header.Get("Content-Type"))
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		if got := r.Form.Get("secret"); got != "private" {
			t.Errorf("secret = %q", got)
		}
		if got := r.Form.Get("response"); got != "token" {
			t.Errorf("response = %q", got)
		}
		if _, sent := r.Form["remoteip"]; sent {
			t.Error("client IP should not be sent to Google")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	a := &app.App{
		AuthPageCaptchaProvider:  app.CaptchaProviderGoogle,
		AuthPageCaptchaSiteKey:   "public",
		AuthPageCaptchaSecretKey: "private",
	}
	provider := googleRecaptchaProvider{verifyURL: server.URL}
	if got := provider.tokenField(); got != "g-recaptcha-response" {
		t.Fatalf("token field = %q", got)
	}
	if !provider.verify(context.Background(), a, "token") {
		t.Fatal("valid provider response was rejected")
	}
	if !captchaConfigured(a) {
		t.Fatal("configured Google provider was not recognized")
	}
	if captchaConfigured(&app.App{AuthPageCaptchaProvider: "unknown"}) {
		t.Fatal("unknown provider must not be treated as configured")
	}
}

func TestGoogleRecaptchaProviderRejectsUnsuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"success":false}`))
	}))
	defer server.Close()

	provider := googleRecaptchaProvider{verifyURL: server.URL}
	a := &app.App{AuthPageCaptchaSecretKey: "private"}
	if provider.verify(context.Background(), a, "token") {
		t.Fatal("unsuccessful provider response was accepted")
	}
}
