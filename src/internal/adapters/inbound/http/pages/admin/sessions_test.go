package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestHumanDeviceName(t *testing.T) {
	for _, tc := range []struct {
		name, userAgent, fallback, want string
	}{
		{
			name:      "iPhone",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
			want:      "iPhone",
		},
		{
			name:      "Chrome on Windows",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
			want:      "Chrome on Windows",
		},
		{name: "stored device fallback", fallback: "tablet", want: "تبلت"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := humanDeviceName(tc.userAgent, tc.fallback); got != tc.want {
				t.Fatalf("humanDeviceName() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSessionDrawerNamesDialogAndCopiesDetailIDs(t *testing.T) {
	var output bytes.Buffer
	session := SessionRow{ID: "session-id", AppID: "app-id", Identifier: "person@example.com"}
	if err := sessionDrawer("session-detail", session).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, want := range []string{
		`<dialog class="drawer" aria-label="جزئیات نشست"`,
		`data-copy-text="session-id"`,
		`data-copy-text="app-id"`,
		`data-copy-text="person@example.com"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("session drawer is missing %q", want)
		}
	}
}
