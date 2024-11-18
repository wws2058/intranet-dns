package middleware

import "github.com/gin-gonic/gin"

// api limiter, prevent burst requests
func ApiLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func getApiLimiter() {

}
