package services

import (
	"context"
	"log"
	"strings"
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/event"
	"aman/internal/ports/outbound"
)

// WebhookService delivers auth events to the app configured webhook URL.
// Delivery is fire-and-forget: a slow or broken webhook must never block
// authentication.
type WebhookService struct {
	sender outbound.WebhookSender
}

func NewWebhookService(sender outbound.WebhookSender) *WebhookService {
	return &WebhookService{sender: sender}
}

// Emit sends an event to the app webhook if one is configured.
func (w *WebhookService) Emit(_ context.Context, a *app.App, name string, data map[string]any) {
	if a == nil || strings.TrimSpace(a.WebhookURL) == "" {
		return
	}

	e := event.Event{
		Name:  name,
		AppID: a.ID.String(),
		At:    time.Now().UTC(),
		Data:  data,
	}
	url := a.WebhookURL

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := w.sender.Send(ctx, url, e); err != nil {
			log.Printf("webhook: %s to %s failed: %v", e.Name, url, err)
		}
	}()
}
