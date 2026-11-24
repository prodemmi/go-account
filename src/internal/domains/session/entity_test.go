package session

import (
	"testing"
	"time"
)

func TestDisplayStatusMatchesExpiryBoundary(t *testing.T) {
	now := time.Now()
	for _, test := range []struct {
		name string
		s    Session
		want string
	}{
		{name: "active", s: Session{Status: StatusActive, ExpiresAt: now.Add(time.Second)}, want: StatusActive},
		{name: "expired active", s: Session{Status: StatusActive, ExpiresAt: now}, want: "expired"},
		{name: "expired pending", s: Session{Status: StatusPending, ExpiresAt: now.Add(-time.Second)}, want: "expired"},
		{name: "revoked remains revoked", s: Session{Status: StatusRevoked, ExpiresAt: now.Add(-time.Second)}, want: StatusRevoked},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := test.s.DisplayStatus(now); got != test.want {
				t.Fatalf("DisplayStatus() = %q, want %q", got, test.want)
			}
		})
	}
}
