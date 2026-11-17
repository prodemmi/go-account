package notification

import (
	"errors"
	"fmt"

	"aman/internal/domains/notification"
	"aman/internal/ports/outbound"
)

var _ outbound.INotificationService = (*SMSService)(nil)

// SMSService sends SMS notifications via a provider. It wraps provider errors
// into domain errors so callers (services/adapters) handle a single error type.
type SMSService struct {
	provider SMSProvider
}

// SMSProvider is the external SMS gateway interface. In production this would
// wrap Twilio, AWS SNS, etc.
type SMSProvider interface {
	Send(to, message string) (messageID string, err error)
}

// NewSMSService constructs an SMSService backed by the given provider.
func NewSMSService(provider SMSProvider) *SMSService {
	return &SMSService{provider: provider}
}

// SendSMS sends an SMS to the given recipient. It returns domain errors
// (notification.ErrSendFailed / notification.ErrInvalid) rather than raw
// provider errors.
func (s *SMSService) SendSMS(to, message string) error {
	if to == "" || message == "" {
		return notification.ErrInvalid
	}

	_, err := s.provider.Send(to, message)
	if err != nil {
		// Wrap the provider error so the caller sees a domain error with the
		// underlying cause preserved.
		return fmt.Errorf("%w: %v", notification.ErrSendFailed, err)
	}
	return nil
}

// MockProvider is a test/dummy provider that can be configured to succeed or
// fail. It is provided via fx in tests.
type MockProvider struct {
	Fail bool
}

func (m *MockProvider) Send(to, message string) (string, error) {
	if m.Fail {
		return "", errors.New("sms gateway unavailable")
	}
	return "msg-123", nil
}
