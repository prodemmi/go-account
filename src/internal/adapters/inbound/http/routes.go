package httpadapter

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"aman/docs"
	"aman/internal/adapters/inbound/http/handlers"
	"aman/internal/adapters/inbound/http/middleware"
	"aman/internal/adapters/inbound/http/pages"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes wires the primary HTTP adapters onto the gin engine.
func RegisterRoutes(
	engine *gin.Engine,
	authHandler *handlers.AuthHandler,
	authPageHandler *handlers.AuthPageHandler,
	sessionHandler *handlers.SessionHandler,
	appMiddleware middleware.AppMiddleware,
	jwtMiddleware middleware.JWTMiddleware,
) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	engine.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, "0.8.0")
	})

	// Static CSS files for branded auth pages and admin panel. They are
	// embedded in the binary so the admin panel is styled regardless of the
	// process working directory.
	css, err := fs.Sub(pages.CSS, "css")
	if err == nil {
		engine.StaticFS("/static/css", http.FS(css))
	}
	fonts, err := fs.Sub(pages.Fonts, "fonts")
	if err == nil {
		engine.StaticFS("/static/fonts", http.FS(fonts))
	}

	// Unknown routes redirect to nginx's 404 page, except the /auth namespace,
	// which returns an actual 404 (including when custom auth is enabled).
	engine.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == handlers.NotFoundPath || isAuthRoutePath(c.Request.URL.Path) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Redirect(http.StatusFound, handlers.NotFoundPath)
	})
	engine.NoMethod(func(c *gin.Context) {
		if c.Request.URL.Path == handlers.NotFoundPath || isAuthRoutePath(c.Request.URL.Path) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Redirect(http.StatusFound, handlers.NotFoundPath)
	})

	// Browser auth flow (CAS branded pages). This is what apps redirect to.
	authPages := engine.Group("/auth").Use(authPageHandler.RequireManagedAuth())
	{
		authPages.GET("/login", authPageHandler.Login)
		authPages.GET("/verify-otp", authPageHandler.VerifyOTPPage)
		authPages.POST("/send-otp", authPageHandler.SendOTP)
		authPages.POST("/verify-otp", authPageHandler.VerifyOTP)
		authPages.GET("/check", authPageHandler.Check)
		authPages.GET("/logout", authPageHandler.Logout)
		authPages.POST("/logout", authPageHandler.Logout)
	}

	v1 := engine.Group("/v1")

	auth := v1.Group("/auth").Use(gin.HandlerFunc(appMiddleware))
	{
		auth.POST("/send-otp", authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/verify-ott", authHandler.VerifyOTT)
		auth.GET("/check", authHandler.Check)
		auth.GET("/check-login", authHandler.CheckLogin)
		auth.GET("/me", gin.HandlerFunc(jwtMiddleware), authHandler.Me)
		auth.POST("/logout", authHandler.Logout)
		auth.DELETE("/logout", authHandler.Logout)
	}

	// A user can list and terminate their own sessions. CAS does not manage
	// users: ownership is proven by the app + identifier of the session token.
	user := v1.Group("/user").Use(gin.HandlerFunc(appMiddleware), gin.HandlerFunc(jwtMiddleware))
	{
		user.GET("/sessions", sessionHandler.Index)
		user.DELETE("/sessions", sessionHandler.DeleteBatch)
		user.DELETE("/sessions/:id", sessionHandler.Delete)
	}

	if os.Getenv("APP_ENV") != "production" {
		RegisterSwaggerRoutes(engine)
	}
}

func isAuthRoutePath(path string) bool {
	return path == "/auth" || strings.HasPrefix(path, "/auth/")
}

// RegisterSwaggerRoutes exposes the generated OpenAPI docs. It is never
// registered in production.
func RegisterSwaggerRoutes(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
