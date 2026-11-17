package otp

import "errors"

var (
	ErrNotFound      = errors.New("otp not found")
	ErrExpired       = errors.New("otp code has expired")
	ErrInvalid       = errors.New("invalid otp")
	ErrAlreadyActive = errors.New("otp already active")
)
