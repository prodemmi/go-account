package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type Header struct {
	Token string `header:"TOKEN" binding:"required"`
	AppID string `header:"APP-ID" binding:"required"`
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &Header{}
		if err := c.BindHeader(header); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, err := jwt.Parse(header.Token, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized.",
				})
			}

			return nil, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You're Unauthorized due to error parsing the JWT.",
			})
			return
		}

		if token.Valid {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token is invalid.",
			})
		}
	}
}
