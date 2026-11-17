package admin

import (
	"go.uber.org/fx"
)

// Module provides the admin panel (separate from /v1/* public API routes).
// RegisterAdminRoutes is called directly from NewServer rather than provided
// through fx, because fx requires every provided function to return at least
// one non-error value.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewAdminHandler),
	)
}
