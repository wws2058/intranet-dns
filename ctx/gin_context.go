package ctx

import "github.com/gin-gonic/gin"

// 设置gin request请求context, 通过gin.Context暂存的在每个请求生命周期内有效的数据

const (
	prefix       = "dns-service_"
	keyRequestID = prefix + "request-id"
)

// 设置request id
func SetRequestID(c *gin.Context, uuid string) {
	c.Set(keyRequestID, uuid)
}

// 获取request id
func GetRequestID(c *gin.Context) string {
	return c.GetString(keyRequestID)
}
