package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"

	pageData "aman/internal/adapters/inbound/http/pages"
	"aman/internal/domains/app"
)

func TestGoogleCaptchaIsRenderedForLoginAndRetry(t *testing.T) {
	data := pageData.AuthData{
		App:             app.App{Name: "Example"},
		AppID:           "app-id",
		CaptchaEnabled:  true,
		CaptchaProvider: app.CaptchaProviderGoogle,
		CaptchaSiteKey:  "public-key",
	}
	var output bytes.Buffer
	if err := Layout(data, authCard(data)).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	if !strings.Contains(html, "https://www.google.com/recaptcha/api.js") {
		t.Fatal("Google reCAPTCHA client script is missing")
	}
	if got := strings.Count(html, `class="g-recaptcha"`); got != 2 {
		t.Fatalf("rendered %d Google widgets, want one for login and one for retry", got)
	}
	if !strings.Contains(html, `name="retry" value="1"`) || strings.Contains(html, `class="arcaptcha"`) {
		t.Fatal("retry form/provider markup does not match the configured Google provider")
	}
}

func TestVerifyPageRendersExpiryAndChannelSpecificEditLabel(t *testing.T) {
	data := pageData.AuthData{
		App:          app.App{Name: "Example"},
		AppID:        "app-id",
		Identifier:   "person@example.com",
		OtpHash:      "signed-otp",
		OTPExpiresAt: 1_800_000_000,
		AllowEmail:   true,
	}
	var output bytes.Buffer
	if err := VerifyOTP(data).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, want := range []string{
		`data-auth-step="verify" data-expires="1800000000"`,
		`name="expires_at" value="1800000000"`,
		`data-auth-back>تغییر ایمیل</a>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("verify page is missing %q", want)
		}
	}
}

func TestEditIdentifierTextMatchesBusinessChannel(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		allowEmail, allowMobile bool
		want                    string
	}{
		{name: "mobile", allowMobile: true, want: "تغییر شماره موبایل"},
		{name: "email", allowEmail: true, want: "تغییر ایمیل"},
		{name: "both", allowEmail: true, allowMobile: true, want: "تغییر ایمیل یا شماره موبایل"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := pageData.EditIdentifierText(tc.allowEmail, tc.allowMobile); got != tc.want {
				t.Fatalf("EditIdentifierText() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCustomAuthLayoutRendersConfiguredColorsAndStructure(t *testing.T) {
	data := pageData.AuthData{
		App:                 app.App{Name: "Example"},
		Layout:              "split",
		TextColor:           "#123456",
		MutedTextColor:      "#654321",
		CardBackgroundColor: "#abcdef",
	}
	var output bytes.Buffer
	if err := Layout(data, authCard(data)).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	html := output.String()
	for _, want := range []string{
		`class="layout-split"`,
		`--auth-text:#123456`,
		`--auth-muted:#654321`,
		`--auth-card:#abcdef`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("custom auth page is missing %q", want)
		}
	}
}
