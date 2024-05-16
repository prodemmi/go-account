package middlewares

import (
	"account/http/services"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWT(userService services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("API-TOKEN")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized.",
			})
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("APP_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token is expired",
			})
			return
		}

		user, errUser := userService.GetUserByToken(token.Raw)

		if errUser != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		ctx.Set("Token", token.Raw)
		ctx.Set("User", user)
		ctx.Set("UserID", user.ID)
		ctx.Next()
	}
}

func NewJWTMiddleware(userService services.IUserService) gin.HandlerFunc {
	return JWT(userService)
}
