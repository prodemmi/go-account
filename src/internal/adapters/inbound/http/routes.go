package httpadapter

import (
	"net/http"

	"aman/docs"
	"aman/internal/adapters/inbound/http/handlers"
	"aman/internal/adapters/inbound/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes wires the primary HTTP adapters onto the gin engine.
func RegisterRoutes(
	engine *gin.Engine,
	authHandler *handlers.AuthHandler,
	appHandler *handlers.AppHandler,
	sessionHandler *handlers.SessionHandler,
	appMiddleware middleware.AppMiddleware,
	jwtMiddleware middleware.JWTMiddleware,
) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	engine.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, "0.7.0")
	})

	// Static CSS files for branded auth pages and admin panel
	engine.StaticFS("/static/css", http.Dir("static/css"))

	v1 := engine.Group("/v1")

	auth := v1.Group("/auth").Use(gin.HandlerFunc(appMiddleware))
	{
		auth.GET("/me", gin.HandlerFunc(jwtMiddleware), authHandler.Me)
		auth.POST("/sign-in", authHandler.SignIn)
		auth.DELETE("/logout", gin.HandlerFunc(jwtMiddleware), authHandler.Logout)
		auth.POST("/verify-otp", authHandler.VerifyOtp)
		auth.POST("/sign-up", authHandler.SignUp)
	}

	apps := v1.Group("/apps").Use(gin.HandlerFunc(appMiddleware))
	{
		apps.POST("/", appHandler.Store)
		apps.GET("/", appHandler.Index)
		apps.GET("/:id", appHandler.Single)
		apps.PUT("/:id", appHandler.Update)
		apps.DELETE("/:id", appHandler.Delete)
	}

	sessions := v1.Group("/sessions").Use(gin.HandlerFunc(appMiddleware))
	{
		sessions.POST("/", sessionHandler.Create)
		sessions.GET("/validate", sessionHandler.Validate)
		sessions.DELETE("/", sessionHandler.Delete)
	}

	RegisterSwaggerRoutes(engine)
}

func RegisterSwaggerRoutes(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
