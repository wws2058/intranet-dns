package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wws2058/intranet-dns/ctx"
	"github.com/wws2058/intranet-dns/models"
)

// @Summary  ping
// @Tags     health
// @Param    mock  query   string           false  "mock"
// @Success  200   object  ctx.StdResponse  "pong"
// @Router   /api/v1/ping [GET]
func pingHandler(c *gin.Context) {
	ctx.SucceedRsp(c, "pong", nil)
}

// register ping test api in router's engine
func LoadPingApis(r *gin.Engine) {
	apis := []models.Api{
		{Path: "/ping", Method: http.MethodGet, Description: "服务测试ping接口", Handler: pingHandler},
	}
	loadApi(r, ginGroupApiV1, apis)
}
