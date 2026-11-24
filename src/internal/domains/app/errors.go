package app

import "errors"

var (
	ErrNotFound          = errors.New("app not found")
	ErrAlreadyExists     = errors.New("app already exists")
	ErrInvalid           = errors.New("invalid app")
	ErrForbidden         = errors.New("app forbidden")
	ErrCountryRestricted = errors.New("country is not allowed")
	ErrInvalidChannel    = errors.New("invalid otp channel")
)
