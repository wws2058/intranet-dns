package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/intranet-dns/ctx"
	"golang.org/x/time/rate"
)

// api limiter, prevent burst requests
func ApiLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rate.NewLimiter(5, 10)
		if !limiter.Allow() {
			ctx.AbortRsp(c, fmt.Errorf("apis rate limit"))
			return
		}
		c.Next()
	}
}
