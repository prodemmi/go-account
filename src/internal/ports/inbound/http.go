// Package inbound contains the contracts exposed to the inbound (HTTP)
// adapters. The concrete services in domains/services satisfy these
// interfaces; the services themselves stay plain structs and inject outbound
// adapters to do the work.
package inbound

import (
	"context"
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/session"
)

// SendOTPCommand asks CAS to create and deliver an OTP code.
type SendOTPCommand struct {
	AppID      string
	Identifier string // mobile number or email address
	UserAgent  string
	IP         string
	Country    string
	City       string
	Region     string
	Device     string
}

// OTPChallenge is returned after an OTP has been created and delivered.
type OTPChallenge struct {
	AppID      string
	Identifier string
	Channel    string
	OTPHash    string
	Length     int
	ExpiresAt  time.Time
}

// VerifyOTPCommand verifies a code and creates a session.
type VerifyOTPCommand struct {
	AppID      string
	Identifier string
	Code       string
	OTPHash    string
	UserAgent  string
	IP         string
	Country    string
	City       string
	Region     string
	Device     string
	// ReturnTo is the URL the browser should be sent back to after login when
	// the app allows returning to its own origin.
	ReturnTo string
}

// AuthResult is the outcome of a successful OTP verification.
type AuthResult struct {
	App         *app.App
	Session     *session.Session
	OTT         string
	RedirectURL string
}

// CheckCommand checks whether an existing CAS session is still valid.
type CheckCommand struct {
	AppID        string
	SessionToken string
	UserAgent    string
	ReturnTo     string
}

// CheckResult tells the caller whether the browser already has a CAS session.
type CheckResult struct {
	Authenticated bool
	OTT           string
	RedirectURL   string
}

// VerifyOTTCommand exchanges a one-time ticket for credentials.
type VerifyOTTCommand struct {
	AppID     string
	OTT       string
	UserAgent string
	IP        string
}

// Credentials are returned to an app after it verifies an OTT.
type Credentials struct {
	AppID        string    `json:"app_id"`
	SessionID    string    `json:"session_id"`
	SessionToken string    `json:"session_token"`
	Identifier   string    `json:"identifier"`
	Mobile       string    `json:"mobile,omitempty"`
	Email        string    `json:"email,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	Infinite     bool      `json:"infinite"`
}

// LogoutCommand terminates the session identified by the token.
type LogoutCommand struct {
	AppID        string
	SessionToken string
}

// AuthUseCase is the inbound port for the authentication flow.
type AuthUseCase interface {
	SendOTP(ctx context.Context, cmd SendOTPCommand) (*OTPChallenge, error)
	VerifyOTP(ctx context.Context, cmd VerifyOTPCommand) (*AuthResult, error)
	Check(ctx context.Context, cmd CheckCommand) (*CheckResult, error)
	VerifyOTT(ctx context.Context, cmd VerifyOTTCommand) (*Credentials, error)
	Logout(ctx context.Context, cmd LogoutCommand) error
}

// AppInput is the editable configuration of an app in the admin panel.
type AppInput struct {
	Name        string
	URL         string
	CallbackURL string
	WebhookURL  string

	OTPChannel       string
	OTPLength        int
	OTPCharset       string
	OTPExpireSeconds int
	OTTExpireSeconds int

	SessionTimeoutSeconds int
	InfiniteSession       bool
	SessionCookieName     string
	SessionCookieHTTPOnly bool
	SessionCookieSecure   bool
	SessionCookieSameSite string
	RedirectBack          bool

	CustomAuthEnabled           bool
	AllowedCountries            string // comma-separated ISO 3166-1 alpha-2 country codes
	AuthPageEnabled             bool
	AuthPageURL                 string
	AuthPageTitle               string
	AuthPageLayout              string
	AuthPageFont                string
	AuthPageTextAlign           string
	AuthPageLogo                string
	AuthPageLogoMask            string
	AuthPageLogoPosition        string
	AuthPagePrimaryColor        string
	AuthPageBackgroundColor     string
	AuthPageCardBackgroundColor string
	AuthPageTextColor           string
	AuthPageMutedTextColor      string
	AuthPagePinBackgroundColor  string
	AuthPagePinTextColor        string
	AuthPageTheme               string
	AuthPageCaptchaEnabled      bool
	AuthPageCaptchaProvider     string
	AuthPageCaptchaSiteKey      string
	AuthPageCaptchaSecretKey    string
	AuthPageBannerMode          string
	AuthPageBanners             []app.Banner
	AuthPageCopyCustom          bool

	AuthPageLoginTitle      string
	AuthPageLoginSubtitle   string
	AuthPageIdentifierLabel string
	AuthPageSendButton      string
	AuthPageVerifyTitle     string
	AuthPageVerifySubtitle  string
	AuthPageVerifyButton    string
	AuthPageRetryText       string
	AuthPageBackText        string
	AuthPageFooterText      string
}

// AppFilter narrows down an app list and paginates it.
type AppFilter struct {
	Search       string
	CreatedFrom  *time.Time
	CreatedUntil *time.Time
	Limit        int
	Offset       int
}

// AppUseCase is the inbound port used by the admin panel to manage apps.
type AppUseCase interface {
	List(ctx context.Context) ([]app.App, error)
	ListPage(ctx context.Context, filter AppFilter) ([]app.App, int64, error)
	Get(ctx context.Context, id string) (*app.App, error)
	Create(ctx context.Context, in AppInput) (app.App, error)
	Update(ctx context.Context, id string, in AppInput) (app.App, error)
	Delete(ctx context.Context, id string) error
}

// SessionFilter narrows down an admin session list and paginates it.
type SessionFilter struct {
	AppID        string
	Search       string
	Status       string
	CreatedFrom  *time.Time
	CreatedUntil *time.Time
	Limit        int
	Offset       int
}

// SessionUseCase is the inbound port used by users to manage their sessions.
type SessionUseCase interface {
	List(ctx context.Context, appID, identifier string) ([]session.Session, error)
	ListAll(ctx context.Context, filter SessionFilter) ([]session.Session, int64, error)
	Stats(ctx context.Context, appID, status string, since, until *time.Time) (SessionStats, error)
	DailyCounts(ctx context.Context, appID, status string, since time.Time, until *time.Time) ([]DailyCount, error)
	Terminate(ctx context.Context, appID, identifier, id string) error
	TerminateMany(ctx context.Context, appID, identifier string, ids []string) error
	TerminateAll(ctx context.Context, appID, identifier string) error
	// Revoke is the admin-only unconditional revocation.
	Revoke(ctx context.Context, id string) error
}

// SessionStats is the dashboard aggregate for all apps or one app.
type SessionStats struct {
	Total    int64
	Active   int64
	Pending  int64
	Expired  int64
	Revoked  int64
	Infinite int64
	Today    int64
	Unique   int64
}

// DailyCount is one point on the dashboard sessions chart.
type DailyCount struct {
	Day   time.Time
	Count int64
}
