package models

import (
	"testing"

	"aman/internal/domains/app"

	"github.com/google/uuid"
)

func TestCustomAuthEnabledRoundTrips(t *testing.T) {
	id := uuid.New()
	model := ToModelApp(app.App{
		ID:                      id,
		CustomAuthEnabled:       true,
		AllowedCountries:        []string{"CA", "US"},
		AuthPageCaptchaProvider: app.CaptchaProviderGoogle,
	})
	if !model.CustomAuthEnabled {
		t.Fatal("custom auth flag was not persisted")
	}
	if got := ToDomainApp(model); got.ID != id || !got.CustomAuthEnabled {
		t.Fatalf("custom auth flag did not load: %+v", got)
	}
	got := ToDomainApp(model)
	if len(got.AllowedCountries) != 2 || got.AllowedCountries[0] != "CA" || got.AllowedCountries[1] != "US" {
		t.Fatalf("allowed countries did not round-trip: %#v", got.AllowedCountries)
	}
	if got.AuthPageCaptchaProvider != app.CaptchaProviderGoogle {
		t.Fatalf("captcha provider did not round-trip: %q", got.AuthPageCaptchaProvider)
	}
}
