package apis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/dns-service/ctx"
	"github.com/tswcbyy1107/dns-service/models"
)

// gin group
var ginGroupApiV1 = "/api/v1"

// standard response
type StdResponse struct {
	Status    bool            `json:"status,omitempty"`     // true: succeed, false: failed
	Pages     *models.PageRsp `json:"pages,omitempty"`      // pages
	Data      interface{}     `json:"data,omitempty"`       // data
	RequestID string          `json:"request_id,omitempty"` // api request uid
	Error     *models.Errors  `json:"error,omitempty"`      // self serr
}

// request succeed, http code 200, status true
func succeedRsp(c *gin.Context, data interface{}, pages *models.PageRsp) {
	uid := ctx.GetRequestID(c)

	c.JSON(http.StatusOK, StdResponse{
		Pages:     pages,
		Status:    true,
		Data:      data,
		RequestID: uid,
	})
}

// request failed, http code 200, status false
func failedRsp(c *gin.Context, err *models.Errors) {
	uid := ctx.GetRequestID(c)

	c.JSON(http.StatusOK, StdResponse{
		Status:    false,
		RequestID: uid,
		Error:     err,
	})
}

// gin load api
func loadApi(r *gin.Engine, group string, apis []models.Api) {
	newApis := []models.Api{}
	// gin router
	if len(group) == 0 {
		group = ginGroupApiV1
	}
	gr := r.Group(group)
	for _, api := range apis {
		gr.Handle(api.Method, api.Path, api.Handler)
		api.Path = fmt.Sprintf("%s%s", group, api.Path)
		api.Active = true
		// audited in Non-get requests
		if api.Method != http.MethodGet {
			api.Audit = true
		}
		newApis = append(newApis, api)
	}
	// db save
	models.ApisRegisterInDB(newApis)
}
