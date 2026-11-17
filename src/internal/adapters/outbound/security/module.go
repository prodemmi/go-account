package security

import (
	"aman/internal/ports/outbound"
	"go.uber.org/fx"
)

// Module provides the security outbound adapters as their port interfaces so
// fx can inject them where the ports are required.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() outbound.IPasswordService {
			return NewPasswordService()
		}),
		fx.Provide(func() outbound.IJWTService {
			return NewJWTService()
		}),
		fx.Provide(func() outbound.IAuthHash {
			return NewAuthHash()
		}),
	)
}
