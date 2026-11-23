package httpadapter

import (
	"aman/internal/adapters/inbound/http/admin"
	"aman/internal/adapters/inbound/http/handlers"
	"aman/internal/adapters/inbound/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module provides the primary HTTP inbound adapter: the gin engine, handlers,
// middlewares, route registration (in NewServer) and the lifecycle-managed
// server.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(gin.Default),
		fx.Provide(handlers.NewAuthHandler),
		fx.Provide(handlers.NewAuthPageHandler),
		fx.Provide(handlers.NewSessionHandler),
		fx.Provide(middleware.NewAppMiddleware),
		fx.Provide(middleware.NewJWTMiddleware),
		fx.Provide(NewServer),
		// fx only constructs providers reachable from an fx.Invoke target.
		// NewServer registers the HTTP OnStart hook, so we invoke it here.
		fx.Invoke(func(s *Server) {}),
		admin.Module(),
	)
}
