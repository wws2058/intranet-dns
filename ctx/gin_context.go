package ctx

import "github.com/gin-gonic/gin"

// set gin context data

const (
	prefix       = "dns-service_"
	keyRequestID = prefix + "request-id"
	keySecretApi = prefix + "sensitive-api"
	keyUserName  = prefix + "user-name"
)

// set request id
func SetRequestID(c *gin.Context, uuid string) {
	c.Set(keyRequestID, uuid)
}

// get request id
func GetRequestID(c *gin.Context) string {
	return c.GetString(keyRequestID)
}

// set sensitive api
func SetSensitiveApi(c *gin.Context) {
	c.Set(keySecretApi, true)
}

// get sensitive api
func GetSensitiveApi(c *gin.Context) bool {
	return c.GetBool(keySecretApi)
}

// set username
func SetLoginUsername(c *gin.Context, name string) {
	c.Set(keyUserName, name)
}

// get username
func GetLoginUsername(c *gin.Context) string {
	return c.GetString(keyUserName)
}
