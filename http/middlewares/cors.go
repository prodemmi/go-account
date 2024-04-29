package middlewares

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Headers", "APP-ID")
		} else {
			// Your code goes here
		}

		c.Next()
	}
}

func NewCORSMiddleware() gin.HandlerFunc {
	return CORS()
}
