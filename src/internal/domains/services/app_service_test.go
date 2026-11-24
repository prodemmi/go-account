package services

import (
	"testing"

	"aman/internal/domains/app"
	"aman/internal/ports/inbound"

	"github.com/google/uuid"
)

func TestAppFromInputNormalizesAuthConfiguration(t *testing.T) {
	got, err := appFromInput(uuid.New(), inbound.AppInput{
		Name:                    "Example",
		CallbackURL:             "https://example.com/callback",
		OTPLength:               5,
		AllowedCountries:        " ca, US,ca ",
		AuthPageCaptchaProvider: app.CaptchaProviderGoogle,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.AllowedCountries) != 2 || got.AllowedCountries[0] != "CA" || got.AllowedCountries[1] != "US" {
		t.Fatalf("countries = %#v, want [CA US]", got.AllowedCountries)
	}
	if got.AuthPageCaptchaProvider != app.CaptchaProviderGoogle {
		t.Fatalf("captcha provider = %q", got.AuthPageCaptchaProvider)
	}
}

func TestAppFromInputRejectsMalformedCountryCode(t *testing.T) {
	for _, raw := range []string{"USA", "U1", "A"} {
		t.Run(raw, func(t *testing.T) {
			_, err := appFromInput(uuid.New(), inbound.AppInput{
				Name:             "Example",
				CallbackURL:      "https://example.com/callback",
				OTPLength:        5,
				AllowedCountries: raw,
			})
			if err == nil {
				t.Fatalf("expected invalid country code %q to be rejected", raw)
			}
		})
	}
}
