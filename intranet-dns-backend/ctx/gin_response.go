package ctx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/intranet-dns/models"
)

// standard response
type StdResponse struct {
	Status    bool            `json:"status"`               // true: succeed, false: failed
	Pages     *models.PageRsp `json:"pages,omitempty"`      // pages
	Data      interface{}     `json:"data,omitempty"`       // data
	RequestID string          `json:"request_id,omitempty"` // api request uid
	Error     string          `json:"error,omitempty"`      // self err
}

// request succeed, http code 200, status true
func SucceedRsp(c *gin.Context, data interface{}, pages *models.PageRsp) {
	uid := GetRequestID(c)

	c.JSON(http.StatusOK, StdResponse{
		Pages:     pages,
		Status:    true,
		Data:      data,
		RequestID: uid,
	})
}

// request failed, http code 200, status false
func FailedRsp(c *gin.Context, err error) {
	uid := GetRequestID(c)

	c.JSON(http.StatusOK, StdResponse{
		Status:    false,
		RequestID: uid,
		Error:     err.Error(),
	})
}

// request failed, http code 200, status false
func AbortRsp(c *gin.Context, err error) {
	uid := GetRequestID(c)

	c.AbortWithStatusJSON(http.StatusOK, StdResponse{
		Status:    false,
		RequestID: uid,
		Error:     err.Error(),
	})
}
