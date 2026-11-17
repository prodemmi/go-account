package otp

import (
	"time"

	"github.com/google/uuid"
)

// Otp represents a one-time password used for app verification in CAS.
type Otp struct {
	ID       uuid.UUID
	AppID    uuid.UUID
	Type     string
	Code     string
	Used     bool
	Hash     *string
	ExpireAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
