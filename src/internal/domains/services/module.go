package services

import (
	"aman/internal/ports/inbound"
	"go.uber.org/fx"
)

// Module provides the application services. The services are concrete
// structs; they are exposed to the HTTP adapters through the inbound ports.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewWebhookService),
		fx.Provide(fx.Annotate(NewAppService, fx.As(new(inbound.AppUseCase)))),
		fx.Provide(fx.Annotate(NewSessionService, fx.As(new(inbound.SessionUseCase)))),
		fx.Provide(fx.Annotate(NewAuthService, fx.As(new(inbound.AuthUseCase)))),
	)
}
