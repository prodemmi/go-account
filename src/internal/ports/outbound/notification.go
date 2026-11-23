package outbound

// Notifier is the outbound port implemented by the notification adapter.
type Notifier interface {
	SendSMS(to, message string) error
	SendEmail(to, subject, body string) error
}
