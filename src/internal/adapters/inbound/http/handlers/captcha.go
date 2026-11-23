package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"aman/internal/domains/app"
)

// ARCaptcha integration. The widget script inserts a hidden input named
// "arcaptcha-token" into the form; the server exchanges that token for a
// success flag at the ARCaptcha verify endpoint.
//
// Keys are stored per app in the admin panel (auth page settings), so every
// app can use its own ARCaptcha account without redeploying the binary.
const arCaptchaVerifyURL = "https://api.arcaptcha.co/arcaptcha/api/verify"

// arCaptchaConfigured reports whether both keys are present. When an app turns
// the captcha on but has no keys yet, the widget is hidden and the login still
// works instead of locking users out.
func arCaptchaConfigured(a *app.App) bool {
	return strings.TrimSpace(a.AuthPageCaptchaSiteKey) != "" &&
		strings.TrimSpace(a.AuthPageCaptchaSecretKey) != ""
}

// verifyArCaptcha checks a challenge token with the ARCaptcha API. Any network
// or decoding failure is treated as a failed challenge.
func verifyArCaptcha(a *app.App, token string) bool {
	token = strings.TrimSpace(token)
	if token == "" || !arCaptchaConfigured(a) {
		return false
	}

	payload, err := json.Marshal(map[string]string{
		"challenge_id": token,
		"site_key":     strings.TrimSpace(a.AuthPageCaptchaSiteKey),
		"secret_key":   strings.TrimSpace(a.AuthPageCaptchaSecretKey),
	})
	if err != nil {
		return false
	}

	req, err := http.NewRequest(http.MethodPost, arCaptchaVerifyURL, bytes.NewReader(payload))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<10))
	if err != nil {
		return false
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false
	}
	return result.Success
}
