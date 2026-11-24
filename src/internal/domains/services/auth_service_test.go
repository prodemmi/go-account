package services

import (
	"context"
	"testing"
	"time"

	"aman/internal/adapters/outbound/notification"
	"aman/internal/adapters/outbound/otpgen"
	"aman/internal/adapters/outbound/postgres"
	"aman/internal/adapters/outbound/security"
	"aman/internal/adapters/outbound/webhook"
	"aman/internal/domains/app"
	"aman/internal/domains/otp"
	"aman/internal/domains/ott"
	"aman/internal/domains/session"
	"aman/internal/ports/inbound"
	"aman/internal/ports/outbound"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOTPToOTTCredentialsFlow(t *testing.T) {
	appID := uuid.New()
	sessionID := uuid.New()
	userAgent := "test-agent"
	identifier := "user@example.com"

	appEntity := app.App{
		ID:               appID,
		Name:             "Test",
		Url:              "https://app.example.com",
		CallbackURL:      "https://app.example.com/auth/callback",
		OTPChannel:       app.OTPChannelEmail,
		OTPExpire:        time.Minute,
		OTTExpire:        30 * time.Second,
		SessionTimeout:   time.Hour,
		AllowedCountries: []string{"CA"},
	}

	var issued ott.OTT
	var createdSession session.Session
	apps := &postgres.MockAppRepository{
		FindByIDFn: func(id string) (*app.App, error) {
			require.Equal(t, appID.String(), id)
			return &appEntity, nil
		},
	}

	sessions := &postgres.MockSessionRepository{
		CreateFn: func(s session.Session) (session.Session, error) {
			s.ID = sessionID
			createdSession = s
			return s, nil
		},
		FindByIDFn: func(id string) (*session.Session, error) {
			require.Equal(t, sessionID.String(), id)
			s := session.Session{
				ID:          sessionID,
				AppID:       appID,
				Token:       appID.String() + ":" + identifier,
				Identifier:  identifier,
				Email:       identifier,
				IsLoggedIn:  true,
				LastLoginAt: ptrTime(time.Now()),
				ExpiresAt:   time.Now().Add(time.Hour),
			}
			return &s, nil
		},
	}

	otps := &postgres.MockOTPRepository{
		ConsumeFn: func(gotAppID, gotIdentifier, code, hash string) (otp.Otp, error) {
			require.Equal(t, appID.String(), gotAppID)
			require.Equal(t, identifier, gotIdentifier)
			require.Equal(t, "12345", code)
			require.NotEmpty(t, hash)
			return otp.Otp{AppID: appID, Identifier: identifier, Code: code}, nil
		},
	}

	otts := &postgres.MockOTTRepository{
		CreateFn: func(o ott.OTT) (ott.OTT, error) {
			o.ID = uuid.New()
			issued = o
			return o, nil
		},
		FindByTokenFn: func(token string) (*ott.OTT, error) {
			require.Equal(t, issued.Token, token)
			return &issued, nil
		},
	}

	notifier := &notification.MockNotifier{}
	sender := &webhook.MockSender{}

	svc := NewAuthService(
		apps,
		sessions,
		otps,
		otts,
		&security.MockAuthHash{},
		&security.MockJWTService{},
		&otpgen.MockGenerator{Code: "12345"},
		notifier,
		NewWebhookService(sender),
		staticGeoIP{location: outbound.IPLocation{Country: "CA", City: "Toronto", Region: "Ontario"}},
	)

	ctx := context.Background()

	challenge, err := svc.SendOTP(ctx, inbound.SendOTPCommand{
		AppID:      appID.String(),
		Identifier: identifier,
		UserAgent:  userAgent,
		IP:         "203.0.113.4",
	})
	require.NoError(t, err)
	require.Equal(t, app.OTPChannelEmail, challenge.Channel)
	require.Len(t, notifier.Emails, 1)

	result, err := svc.VerifyOTP(ctx, inbound.VerifyOTPCommand{
		AppID:      appID.String(),
		Identifier: identifier,
		Code:       "12345",
		OTPHash:    challenge.OTPHash,
		UserAgent:  userAgent,
		IP:         "203.0.113.4",
	})
	require.NoError(t, err)
	require.Equal(t, "https://app.example.com/auth/callback?ott="+result.OTT, result.RedirectURL)
	require.Equal(t, "CA", createdSession.Country)
	require.Equal(t, "Toronto", createdSession.City)
	require.Equal(t, "Ontario", createdSession.Region)

	creds, err := svc.VerifyOTT(ctx, inbound.VerifyOTTCommand{
		AppID: appID.String(),
		OTT:   result.OTT,
	})
	require.NoError(t, err)
	require.Equal(t, identifier, creds.Identifier)
	require.Equal(t, sessionID.String(), creds.SessionID)
	require.Equal(t, appID.String()+":"+identifier, creds.SessionToken)
}

func TestSendOTPRejectsIdentifierForWrongChannel(t *testing.T) {
	appID := uuid.New()
	apps := &postgres.MockAppRepository{
		FindByIDFn: func(string) (*app.App, error) {
			return &app.App{ID: appID, OTPChannel: app.OTPChannelSMS}, nil
		},
	}

	svc := NewAuthService(
		apps,
		&postgres.MockSessionRepository{},
		&postgres.MockOTPRepository{},
		&postgres.MockOTTRepository{},
		&security.MockAuthHash{},
		&security.MockJWTService{},
		&otpgen.MockGenerator{},
		&notification.MockNotifier{},
		NewWebhookService(&webhook.MockSender{}),
		nil,
	)

	_, err := svc.SendOTP(context.Background(), inbound.SendOTPCommand{
		AppID:      appID.String(),
		Identifier: "user@example.com",
	})
	require.Error(t, err)
}

func TestSendOTPRejectsUnresolvedCountryForAllowlistedApp(t *testing.T) {
	appID := uuid.New()
	apps := &postgres.MockAppRepository{
		FindByIDFn: func(string) (*app.App, error) {
			return &app.App{ID: appID, AllowedCountries: []string{"CA"}}, nil
		},
	}
	notifier := &notification.MockNotifier{}
	svc := NewAuthService(
		apps,
		&postgres.MockSessionRepository{},
		&postgres.MockOTPRepository{},
		&postgres.MockOTTRepository{},
		&security.MockAuthHash{},
		&security.MockJWTService{},
		&otpgen.MockGenerator{Code: "12345"},
		notifier,
		NewWebhookService(&webhook.MockSender{}),
		nil,
	)

	_, err := svc.SendOTP(context.Background(), inbound.SendOTPCommand{
		AppID:      appID.String(),
		Identifier: "user@example.com",
		Country:    "CA", // Request headers are not authoritative for geo-fencing.
		IP:         "203.0.113.4",
	})
	require.ErrorIs(t, err, app.ErrCountryRestricted)
	require.Empty(t, notifier.Emails)
}

type staticGeoIP struct {
	location outbound.IPLocation
	err      error
}

func (s staticGeoIP) Lookup(string) (outbound.IPLocation, error) {
	return s.location, s.err
}

func ptrTime(t time.Time) *time.Time { return &t }
