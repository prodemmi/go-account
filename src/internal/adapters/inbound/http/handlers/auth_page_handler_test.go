package handlers

import "testing"

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
