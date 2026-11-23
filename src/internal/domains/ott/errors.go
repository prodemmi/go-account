package ott

import "errors"

var (
	ErrNotFound  = errors.New("ott not found")
	ErrExpired   = errors.New("ott has expired")
	ErrUsed      = errors.New("ott already used")
	ErrInvalid   = errors.New("invalid ott")
	ErrForbidden = errors.New("ott forbidden")
)
