package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"aman/internal/domains/app"
)

const (
	arCaptchaVerifyURL = "https://api.arcaptcha.co/arcaptcha/api/verify"
	googleVerifyURL    = "https://www.google.com/recaptcha/api/siteverify"
)

type captchaProvider interface {
	siteKey(*app.App) string
	secretKey(*app.App) string
	tokenField() string
	verify(context.Context, *app.App, string) bool
}

type arCaptchaProvider struct{ verifyURL string }
type googleRecaptchaProvider struct{ verifyURL string }

var captchaHTTPClient = &http.Client{Timeout: 10 * time.Second}

func captchaProviderFor(a *app.App) captchaProvider {
	if a == nil {
		return nil
	}
	switch a.AuthPageCaptchaProvider {
	case app.CaptchaProviderARCaptcha:
		return arCaptchaProvider{verifyURL: arCaptchaVerifyURL}
	case app.CaptchaProviderGoogle:
		return googleRecaptchaProvider{verifyURL: googleVerifyURL}
	default:
		return nil
	}
}

func captchaConfigured(a *app.App) bool {
	provider := captchaProviderFor(a)
	return provider != nil && strings.TrimSpace(provider.siteKey(a)) != "" && strings.TrimSpace(provider.secretKey(a)) != ""
}

func verifyCaptcha(ctx context.Context, a *app.App, token string) bool {
	provider := captchaProviderFor(a)
	token = strings.TrimSpace(token)
	return provider != nil && token != "" && captchaConfigured(a) && provider.verify(ctx, a, token)
}

func (arCaptchaProvider) siteKey(a *app.App) string   { return a.AuthPageCaptchaSiteKey }
func (arCaptchaProvider) secretKey(a *app.App) string { return a.AuthPageCaptchaSecretKey }
func (arCaptchaProvider) tokenField() string          { return "arcaptcha-token" }

func (p arCaptchaProvider) verify(ctx context.Context, a *app.App, token string) bool {
	payload, err := json.Marshal(map[string]string{
		"challenge_id": token,
		"site_key":     strings.TrimSpace(a.AuthPageCaptchaSiteKey),
		"secret_key":   strings.TrimSpace(a.AuthPageCaptchaSecretKey),
	})
	if err != nil {
		return false
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.verifyURL, bytes.NewReader(payload))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := captchaHTTPClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false
	}
	return captchaResponseSucceeded(resp.Body)
}

func (googleRecaptchaProvider) siteKey(a *app.App) string   { return a.AuthPageCaptchaSiteKey }
func (googleRecaptchaProvider) secretKey(a *app.App) string { return a.AuthPageCaptchaSecretKey }
func (googleRecaptchaProvider) tokenField() string          { return "g-recaptcha-response" }

func (p googleRecaptchaProvider) verify(ctx context.Context, a *app.App, token string) bool {
	form := url.Values{
		"secret":   {strings.TrimSpace(a.AuthPageCaptchaSecretKey)},
		"response": {token},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := captchaHTTPClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false
	}
	return captchaResponseSucceeded(resp.Body)
}

func captchaResponseSucceeded(body io.Reader) bool {
	response, err := io.ReadAll(io.LimitReader(body, 16<<10))
	if err != nil {
		return false
	}
	var result struct {
		Success bool `json:"success"`
	}
	return json.Unmarshal(response, &result) == nil && result.Success
}
