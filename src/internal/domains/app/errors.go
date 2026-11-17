package app

import "errors"

var (
	ErrNotFound      = errors.New("app not found")
	ErrAlreadyExists = errors.New("app already exists")
	ErrInvalid       = errors.New("invalid app")
)
