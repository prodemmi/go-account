package middlewares

import (
	"account/repositories"
	"account/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
)

func JWT(userService services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("API-TOKEN")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized.",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("APP_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}

		user, errUser := userService.GetUserByToken(token.Raw)

		if errUser != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}

		ctx.Set("Token", token.Raw)
		ctx.Set("User", user)
		ctx.Set("UserID", user.ID)
		ctx.Next()
	}
}

func NewJWTMiddleware() gin.HandlerFunc {
	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(userRepository)

	return JWT(userService)
}
