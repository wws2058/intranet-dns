package middleware

import "github.com/gin-gonic/gin"

// auth
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
