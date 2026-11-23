package session

import "errors"

var (
	ErrNotFound      = errors.New("session not found")
	ErrExpired       = errors.New("session expired")
	ErrInvalid       = errors.New("invalid session")
	ErrForbidden     = errors.New("session forbidden")
	ErrAlreadyExists = errors.New("session already exists")
)
