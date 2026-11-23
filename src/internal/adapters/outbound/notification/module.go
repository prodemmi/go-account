package notification

import (
	"aman/internal/ports/outbound"
	"go.uber.org/fx"
)

// Module provides the notification outbound adapter as its port interface.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() outbound.Notifier {
			return NewService(NewSMSProviderFromEnv(), NewEmailProviderFromEnv())
		}),
	)
}
