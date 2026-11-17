package session

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an authenticated CAS session bound to an app.
type Session struct {
	ID        uuid.UUID
	AppID     uuid.UUID
	Token     string // service ticket / session token
	ExpiresAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
