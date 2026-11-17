package middleware

import (
	"go.uber.org/fx"
)

// Module provides the HTTP middleware providers.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewAppMiddleware),
		fx.Provide(NewJWTMiddleware),
	)
}
