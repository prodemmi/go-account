package security

import (
	"testing"
	"time"

	"aman/internal/domains/app"

	"github.com/google/uuid"
)

func TestAuthHashBindsOTPChallenge(t *testing.T) {
	t.Setenv("APP_SECRET", "test-secret")
	a := app.App{ID: uuid.New()}
	expiresAt := time.Now().Add(time.Minute)
	signer := NewAuthHash()
	token := signer.CreateAuthHash(a, "user@example.com", "browser", expiresAt)

	if !signer.VerifyAuthHash(token, a, "user@example.com", "browser") {
		t.Fatal("valid challenge token was rejected")
	}
	for _, test := range []struct {
		name       string
		app        app.App
		identifier string
		agent      string
		token      string
	}{
		{name: "other app", app: app.App{ID: uuid.New()}, identifier: "user@example.com", agent: "browser", token: token},
		{name: "other identifier", app: a, identifier: "other@example.com", agent: "browser", token: token},
		{name: "other user agent", app: a, identifier: "user@example.com", agent: "other browser", token: token},
		{name: "tampered", app: a, identifier: "user@example.com", agent: "browser", token: "x" + token},
	} {
		t.Run(test.name, func(t *testing.T) {
			if signer.VerifyAuthHash(test.token, test.app, test.identifier, test.agent) {
				t.Fatal("invalid challenge token was accepted")
			}
		})
	}

	expired := signer.CreateAuthHash(a, "user@example.com", "browser", time.Now().Add(-time.Second))
	if signer.VerifyAuthHash(expired, a, "user@example.com", "browser") {
		t.Fatal("expired challenge token was accepted")
	}
	t.Setenv("APP_SECRET", "")
	if signer.CreateAuthHash(a, "user@example.com", "browser", expiresAt) != "" || signer.VerifyAuthHash(token, a, "user@example.com", "browser") {
		t.Fatal("challenge signing did not fail closed without APP_SECRET")
	}
}
