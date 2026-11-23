package outbound

import (
	"context"

	"aman/internal/domains/event"
)

// WebhookSender is the outbound port implemented by the webhook adapter.
type WebhookSender interface {
	Send(ctx context.Context, url string, e event.Event) error
}
