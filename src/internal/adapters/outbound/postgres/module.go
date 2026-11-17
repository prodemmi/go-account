package postgres

import (
	"aman/internal/adapters/outbound/postgres/repos"
	"aman/internal/ports/outbound"
	"gorm.io/gorm"
	"go.uber.org/fx"
)

// Module provides the persistence outbound adapter: the gorm connection and
// the PostgreSQL repositories, exposed as their port interfaces.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewDB),
		fx.Provide(func(db *gorm.DB) outbound.IAppRepository {
			return repos.NewAppRepository(db)
		}),
		fx.Provide(func(db *gorm.DB) outbound.ISessionRepository {
			return repos.NewSessionRepository(db)
		}),
		fx.Provide(func(db *gorm.DB) outbound.IOTPRepository {
			return repos.NewOTPRepository(db)
		}),
	)
}
