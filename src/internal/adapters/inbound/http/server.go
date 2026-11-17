package httpadapter

import (
	"context"
	"log"
	"net/http"
	"os"

	"aman/internal/adapters/inbound/http/admin"
	"aman/internal/adapters/inbound/http/handlers"
	"aman/internal/adapters/inbound/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Server wraps a gin engine behind a start/shutdown lifecycle so it can be
// managed by fx (OnStart / OnStop hooks).
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
}

// NewServer builds the HTTP server. Route registration happens here during
// graph construction (before Start) so routes are guaranteed to be wired when
// the server begins listening.
func NewServer(
	lc fx.Lifecycle,
	engine *gin.Engine,
	authHandler *handlers.AuthHandler,
	appHandler *handlers.AppHandler,
	sessionHandler *handlers.SessionHandler,
	authPageHandler *handlers.AuthPageHandler,
	appMiddleware middleware.AppMiddleware,
	jwtMiddleware middleware.JWTMiddleware,
	adminHandler *admin.AdminHandler,
) *Server {
	// Public API routes (/v1/*)
	RegisterRoutes(engine, authHandler, appHandler, sessionHandler, appMiddleware, jwtMiddleware)

	// Auth pages (CAS serves branded login/signup/verify-OTP pages)
	authPages := engine.Group("/auth")
	{
		authPages.GET("/login", authPageHandler.Login)
		authPages.GET("/signup", authPageHandler.Signup)
		authPages.GET("/verify-otp", authPageHandler.VerifyOTP)
	}

	// Admin panel routes (/admin/*)
	admin.RegisterAdminRoutes(engine, adminHandler)

	s := &Server{engine: engine}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				port := os.Getenv("APP_PORT")
				if err := s.Run(":" + port); err != nil {
					log.Fatalf("http server stopped: %v", err)
				}
			}()
			return nil
		},
		OnStop: s.Shutdown,
	})
	return s
}

// Run blocks listening on addr. Intended to be run in a goroutine from the
// fx OnStart hook.
func (s *Server) Run(addr string) error {
	s.httpServer = &http.Server{Addr: addr, Handler: s.engine}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
