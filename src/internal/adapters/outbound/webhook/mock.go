package webhook

import (
	"context"
	"errors"

	"aman/internal/domains/event"
)

// MockSender records events instead of sending them.
type MockSender struct {
	Fail   bool
	Events []event.Event
	URLs   []string
}

func (m *MockSender) Send(_ context.Context, url string, e event.Event) error {
	if m.Fail {
		return errors.New("webhook unavailable")
	}
	m.URLs = append(m.URLs, url)
	m.Events = append(m.Events, e)
	return nil
}
