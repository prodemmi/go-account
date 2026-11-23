package notification

import "errors"

// MockNotifier is a test double that records every delivery and can be told
// to fail.
type MockNotifier struct {
	Fail   bool
	SMS    []Message
	Emails []Message
}

type Message struct {
	To      string
	Subject string
	Body    string
}

func (m *MockNotifier) SendSMS(to, message string) error {
	if m.Fail {
		return errors.New("sms gateway unavailable")
	}
	m.SMS = append(m.SMS, Message{To: to, Body: message})
	return nil
}

func (m *MockNotifier) SendEmail(to, subject, body string) error {
	if m.Fail {
		return errors.New("smtp unavailable")
	}
	m.Emails = append(m.Emails, Message{To: to, Subject: subject, Body: body})
	return nil
}
