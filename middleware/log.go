package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tswcbyy1107/dns-service/ctx"
	"github.com/tswcbyy1107/dns-service/utils"
)

// log, set request id, trace api request
func LogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := utils.GenUUID()
		ctx.SetRequestID(c, uid)

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		requestBody := []byte("ignore request body")
		if utils.Contains([]string{"PUT", "POST", "DELETE"}, c.Request.Method) && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		start := time.Now()
		c.Next()
		cost := time.Since(start)

		requestBodyString := string(requestBody)
		if len(requestBodyString) > 1024 {
			requestBodyString = requestBodyString[:1024] + "..."
		}
		requestBodyString = strings.Replace(requestBodyString, "\n", "", -1)
		requestBodyString = strings.Replace(requestBodyString, "\t", "", -1)
		requestBodyString = strings.Replace(requestBodyString, "\r", "", -1)
		requestBodyString = strings.Replace(requestBodyString, " ", "", -1)

		responseBody := "ignore response body"
		if utils.Contains([]string{"PUT", "POST", "DELETE"}, c.Request.Method) {
			responseBody = w.body.String()
			if len(responseBody) > 1024 {
				responseBody = responseBody[:1024] + "..."
			}
		}
		responseBody = strings.Replace(responseBody, "\n", "", -1)
		responseBody = strings.Replace(responseBody, "\t", "", -1)
		responseBody = strings.Replace(responseBody, "\r", "", -1)
		responseBody = strings.Replace(responseBody, " ", "", -1)

		fields := map[string]interface{}{
			"headers":       c.Request.Header,
			"method":        c.Request.Method,
			"url":           c.Request.URL.String(),
			"cost":          cost.Nanoseconds() / 1000000,
			"status_code":   c.Writer.Status(),
			"remote_ip":     c.ClientIP(),
			"request_body":  requestBodyString,
			"response_body": responseBody,
			"request_id":    uid,
		}
		go func() {
			logrus.WithFields(fields).Info("api_detail")
		}()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
