package notification

import "errors"

var (
	ErrSendFailed = errors.New("notification send failed")
	ErrInvalid    = errors.New("invalid notification")
)
