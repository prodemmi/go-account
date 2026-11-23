package session

import (
	"time"

	"github.com/google/uuid"
)

// Session lifecycle statuses. A session created after a successful OTP
// verification waits for the app to exchange its OTT, so it starts as pending
// and only becomes active once VerifyOTT succeeds.
const (
	StatusPending = "pending"
	StatusActive  = "active"
	StatusRevoked = "revoked"
)

// Session represents an authenticated CAS session bound to an app and to a
// single identifier (mobile number or email address). CAS does not manage
// users; the identifier is the only identity it knows about.
type Session struct {
	ID         uuid.UUID `json:"id"`
	AppID      uuid.UUID `json:"app_id"`
	Token      string    `json:"-"` // service ticket / session token
	Identifier string    `json:"identifier"`
	Mobile     string    `json:"mobile,omitempty"`
	Email      string    `json:"email,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	IP         string    `json:"ip,omitempty"`
	Country    string    `json:"country,omitempty"`
	City       string    `json:"city,omitempty"`
	Region     string    `json:"region,omitempty"`
	Device     string    `json:"device,omitempty"`
	Status     string    `json:"status"`
	IsLoggedIn bool      `json:"is_logged_in"`
	Infinite   bool      `json:"infinite"`

	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NormalizedStatus returns the stored status, upgrading legacy rows that only
// have the old is_logged_in flag.
func (s Session) NormalizedStatus() string {
	if s.Status != "" {
		return s.Status
	}
	if s.IsLoggedIn {
		return StatusActive
	}
	return StatusRevoked
}

// Pending reports whether the session is waiting for its OTT to be exchanged.
func (s Session) Pending() bool {
	return s.NormalizedStatus() == StatusPending
}

// Active reports whether the session is still usable at the given time.
func (s Session) Active(now time.Time) bool {
	return s.NormalizedStatus() == StatusActive && s.ExpiresAt.After(now)
}

// DisplayStatus returns the status shown to users, including the derived
// "expired" state for any non-revoked session whose expiry has passed. It is
// used by the admin list and must match the SQL filters in the repository.
func (s Session) DisplayStatus(now time.Time) string {
	status := s.NormalizedStatus()
	if status == StatusRevoked {
		return status
	}
	if !s.ExpiresAt.After(now) {
		return "expired"
	}
	return status
}

// MatchesAgent validates the user agent stored at OTP verification against the
// agent presented on a later step. Empty agents are treated as matching so old
// rows keep working.
func (s Session) MatchesAgent(agent string) bool {
	if s.UserAgent == "" || agent == "" {
		return true
	}
	return s.UserAgent == agent
}
