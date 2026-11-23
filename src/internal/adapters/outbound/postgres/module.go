package postgres

import (
	"aman/internal/adapters/outbound/postgres/repos"
	"aman/internal/ports/outbound"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module provides the persistence outbound adapter: the gorm connection and
// the PostgreSQL repositories, exposed as their port interfaces.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewDB),
		fx.Provide(func(db *gorm.DB) outbound.AppRepository {
			return repos.NewAppRepository(db)
		}),
		fx.Provide(func(db *gorm.DB) outbound.SessionRepository {
			return repos.NewSessionRepository(db)
		}),
		fx.Provide(func(db *gorm.DB) outbound.OTPRepository {
			return repos.NewOTPRepository(db)
		}),
		fx.Provide(func(db *gorm.DB) outbound.OTTRepository {
			return repos.NewOTTRepository(db)
		}),
	)
}
