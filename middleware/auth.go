package middleware

import (
	"github.com/gin-gonic/gin"
)

// TODO
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
