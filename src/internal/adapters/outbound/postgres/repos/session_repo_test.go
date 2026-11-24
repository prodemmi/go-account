package repos

import (
	"reflect"
	"testing"
	"time"

	"aman/internal/domains/session"
)

func TestSessionStatusFilterUsesSessionTableLifecycleRules(t *testing.T) {
	now := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		status string
		query  string
		args   []any
	}{
		{session.StatusActive, legacyActive + " AND expires_at > ?", []any{now}},
		{"expired", expired, []any{now}},
		{session.StatusPending, "status = ? AND expires_at > ?", []any{session.StatusPending, now}},
		{session.StatusRevoked, legacyRevoked, nil},
		{"unknown", "", nil},
	}
	for _, tc := range tests {
		t.Run(tc.status, func(t *testing.T) {
			query, args := sessionStatusFilter(tc.status, now)
			if query != tc.query || !reflect.DeepEqual(args, tc.args) {
				t.Fatalf("sessionStatusFilter(%q) = (%q, %v), want (%q, %v)", tc.status, query, args, tc.query, tc.args)
			}
		})
	}
}
