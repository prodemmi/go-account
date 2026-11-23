package ott

import (
	"time"

	"github.com/google/uuid"
)

// OTT (one-time ticket) is a short-lived opaque ticket handed to the app in
// the callback URL. The app backend exchanges it exactly once for credentials
// through /verify-ott.
type OTT struct {
	ID         uuid.UUID
	AppID      uuid.UUID
	SessionID  uuid.UUID
	Token      string
	Identifier string
	Used       bool
	ExpireAt   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
