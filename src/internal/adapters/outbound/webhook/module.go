package webhook

import (
	"aman/internal/ports/outbound"
	"go.uber.org/fx"
)

// Module provides the webhook outbound adapter as its port interface.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() outbound.WebhookSender {
			return NewClient()
		}),
	)
}
