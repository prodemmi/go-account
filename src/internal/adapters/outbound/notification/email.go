package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"aman/internal/domains/notification"
)

// EmailProvider is the external email gateway interface.
type EmailProvider interface {
	Send(to, subject, body string) (messageID string, err error)
}

// LogEmailProvider prints the message instead of sending it.
type LogEmailProvider struct{}

func (p *LogEmailProvider) Send(to, subject, body string) (string, error) {
	log.Printf("notification/email: to=%s subject=%s body=%s", to, subject, body)
	return "log", nil
}

// SMTPEmailProvider delivers mail through an SMTP server.
type SMTPEmailProvider struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (p *SMTPEmailProvider) Send(to, subject, body string) (string, error) {
	addr := p.Host + ":" + p.Port
	auth := smtp.PlainAuth("", p.Username, p.Password, p.Host)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + p.From + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" + body + "\r\n")

	if err := smtp.SendMail(addr, auth, p.From, []string{to}, msg); err != nil {
		return "", fmt.Errorf("smtp send: %w", err)
	}
	return "smtp", nil
}

// NewEmailProviderFromEnv uses SMTP when SMTP_HOST is configured. Logging is a
// development-only fallback; production refuses to leak OTP codes into logs.
func NewEmailProviderFromEnv() EmailProvider {
	if host := os.Getenv("SMTP_HOST"); host != "" {
		port := os.Getenv("SMTP_PORT")
		if port == "" {
			port = "587"
		}
		from := os.Getenv("SMTP_FROM")
		if from == "" {
			from = os.Getenv("SMTP_USERNAME")
		}
		return &SMTPEmailProvider{
			Host:     host,
			Port:     port,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     from,
		}
	}
	if os.Getenv("APP_ENV") == "development" {
		return &LogEmailProvider{}
	}
	log.Println("notification/email: SMTP_HOST is not set; email delivery is disabled")
	return &unavailableEmailProvider{}
}

// unavailableEmailProvider refuses to send in production when SMTP is not
// configured.
type unavailableEmailProvider struct{}

func (p *unavailableEmailProvider) Send(string, string, string) (string, error) {
	return "", notification.ErrSendFailed
}
