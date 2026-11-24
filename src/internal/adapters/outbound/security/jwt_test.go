package security

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestParseTokenValidatesHMACAndClaims(t *testing.T) {
	t.Setenv("APP_SECRET", "test-secret")

	service := NewJWTService()
	token, err := service.CreateToken("app", "user", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	appID, identifier, err := service.ParseToken(token)
	if err != nil || appID != "app" || identifier != "user" {
		t.Fatalf("ParseToken() = (%q, %q, %v), want (app, user, nil)", appID, identifier, err)
	}

	unsigned := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`)) + "." +
		base64.RawURLEncoding.EncodeToString([]byte(`{"iss":"app","sub":"user"}`)) + "."
	if _, _, err := service.ParseToken(unsigned); err == nil {
		t.Fatal("ParseToken accepted a token using the none algorithm")
	}
}
