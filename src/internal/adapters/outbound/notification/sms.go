package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"aman/internal/domains/notification"
	"aman/internal/ports/outbound"
)

var _ outbound.Notifier = (*Service)(nil)

// Service delivers OTP codes by SMS and email. Provider errors are wrapped
// into domain errors so callers handle a single error type.
type Service struct {
	sms   SMSProvider
	email EmailProvider
}

func NewService(sms SMSProvider, email EmailProvider) *Service {
	return &Service{sms: sms, email: email}
}

func (s *Service) SendSMS(to, message string) error {
	if to == "" || message == "" {
		return notification.ErrInvalid
	}
	if _, err := s.sms.Send(to, message); err != nil {
		return fmt.Errorf("%w: %v", notification.ErrSendFailed, err)
	}
	return nil
}

func (s *Service) SendEmail(to, subject, body string) error {
	if to == "" || subject == "" {
		return notification.ErrInvalid
	}
	if _, err := s.email.Send(to, subject, body); err != nil {
		return fmt.Errorf("%w: %v", notification.ErrSendFailed, err)
	}
	return nil
}

// SMSProvider is the external SMS gateway interface.
type SMSProvider interface {
	Send(to, message string) (messageID string, err error)
}

// LogSMSProvider prints the message instead of sending it. It is the default
// in development so OTP codes are easy to grab from the log.
type LogSMSProvider struct{}

func (p *LogSMSProvider) Send(to, message string) (string, error) {
	log.Printf("notification/sms: to=%s message=%s", to, message)
	return "log", nil
}

// GatewaySMSProvider posts messages to a configurable JSON gateway.
type GatewaySMSProvider struct {
	URL    string
	Token  string
	Client *http.Client
}

func (p *GatewaySMSProvider) Send(to, message string) (string, error) {
	body, err := json.Marshal(map[string]string{"to": to, "message": message})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, p.URL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.Token)
	}

	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("sms gateway returned %s", resp.Status)
	}
	return "gateway", nil
}

// NewSMSProviderFromEnv picks a real gateway when SMS_GATEWAY_URL is set. It
// only falls back to logging in development; production without a gateway
// fails loudly instead of silently logging OTP codes.
func NewSMSProviderFromEnv() SMSProvider {
	if url := os.Getenv("SMS_GATEWAY_URL"); url != "" {
		return &GatewaySMSProvider{URL: url, Token: os.Getenv("SMS_GATEWAY_TOKEN")}
	}
	if os.Getenv("APP_ENV") == "development" {
		return &LogSMSProvider{}
	}
	log.Println("notification/sms: SMS_GATEWAY_URL is not set; SMS delivery is disabled")
	return &unavailableSMSProvider{}
}

// unavailableSMSProvider refuses to send in production when no gateway is
// configured. This prevents OTP codes from leaking into production logs.
type unavailableSMSProvider struct{}

func (p *unavailableSMSProvider) Send(string, string) (string, error) {
	return "", notification.ErrSendFailed
}
