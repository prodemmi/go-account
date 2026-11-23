package otp

import (
	"time"

	"github.com/google/uuid"
)

// Otp represents a one-time password used to authenticate a mobile number or
// email address for an app in CAS.
type Otp struct {
	ID         uuid.UUID
	AppID      uuid.UUID
	Identifier string // mobile number or email address
	Type       string // sms | email
	Code       string
	Used       bool
	Hash       *string
	ExpireAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
