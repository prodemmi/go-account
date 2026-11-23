package event

import "time"

// Event names emitted to an app's webhook URL after every auth action.
const (
	OTPSent        = "otp.sent"
	OTPVerified    = "otp.verified"
	SessionCreated = "session.created"
	SessionRevoked = "session.revoked"
	OTTIssued      = "ott.issued"
	OTTVerified    = "ott.verified"
	LoggedOut      = "auth.logged_out"
)

// Event is the payload POSTed to an app's webhook URL. Apps can use it to
// react to CAS activity (audit trails, notifications, analytics, ...).
type Event struct {
	Name  string         `json:"name"`
	AppID string         `json:"app_id"`
	At    time.Time      `json:"at"`
	Data  map[string]any `json:"data,omitempty"`
}
