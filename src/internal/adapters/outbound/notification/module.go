package notification

import (
	"aman/internal/ports/outbound"
	"go.uber.org/fx"
)

// Module provides the SMS notification outbound adapter as its port interface.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() outbound.INotificationService {
			return NewSMSService(&MockProvider{})
		}),
	)
}
