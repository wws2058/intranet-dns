package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wws2058/intranet-dns/ctx"
	"github.com/wws2058/intranet-dns/models"
	"github.com/wws2058/intranet-dns/utils"
)

// log, set request id, trace api request
func LogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := utils.GenUUID()
		ctx.SetRequestID(c, uid)

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		requestBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		requestBodyString := string(requestBody)
		requestBodyString = strings.Replace(requestBodyString, "\n", "", -1)
		requestBodyString = strings.Replace(requestBodyString, "\t", "", -1)
		requestBodyString = strings.Replace(requestBodyString, "\r", "", -1)

		start := time.Now()
		c.Next()
		cost := time.Since(start)

		responseBodyString := w.body.String()
		responseBodyString = strings.Replace(responseBodyString, "\n", "", -1)
		responseBodyString = strings.Replace(responseBodyString, "\t", "", -1)
		responseBodyString = strings.Replace(responseBodyString, "\r", "", -1)

		// get sensitive mark, no log, no audit
		if ctx.GetSensitiveApi(c) {
			requestBodyString = "{}"
			responseBodyString = "{}"
		}

		fields := map[string]interface{}{
			"method":        c.Request.Method,
			"url":           c.Request.URL.String(),
			"time_cost":     cost.Milliseconds(),
			"status_code":   c.Writer.Status(),
			"remote_ip":     c.ClientIP(),
			"request_id":    uid,
			"request_body":  requestBodyString,
			"response_body": responseBodyString,
			"headers":       c.Request.Header,
		}
		// log and save audit log in db
		go func() {
			if strings.Contains(c.Request.URL.String(), "swagger") {
				return
			}

			logrus.WithFields(fields).Info("api_detail")
			api := &models.Api{}
			dao := &models.DaoDBReq{
				Dst: api,
				ModelFilter: map[string]interface{}{
					"path":   c.FullPath(),
					"method": c.Request.Method,
				},
			}
			err := models.TemplateQuery(dao)
			if err != nil {
				return
			}

			if !api.Audit {
				return
			}

			username := ctx.GetLoginUsername(c)
			models.TemplateCreate(&models.AuditLog{
				UserName:     username,
				RequestID:    uid,
				ClientIP:     c.ClientIP(),
				URL:          c.Request.URL.String(),
				Method:       c.Request.Method,
				RequestBody:  requestBodyString,
				ResponseBody: responseBodyString,
				TimeCost:     int(cost.Milliseconds()),
			})
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
