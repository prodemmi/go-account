package otpgen

import (
	"aman/internal/ports/outbound"
	"go.uber.org/fx"
)

// Module provides the OTP generation outbound adapter as its port interface.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() outbound.OTPGenerator {
			return NewOTPGenerator()
		}),
	)
}
