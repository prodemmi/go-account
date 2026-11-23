package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"aman/internal/domains/event"
	"aman/internal/ports/outbound"
)

var _ outbound.WebhookSender = (*Client)(nil)

// Client delivers events to app webhook URLs over HTTP.
type Client struct {
	HTTP *http.Client
}

func NewClient() *Client {
	return &Client{HTTP: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) Send(ctx context.Context, url string, e event.Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CAS-Webhook/1.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned %s", resp.Status)
	}
	return nil
}
