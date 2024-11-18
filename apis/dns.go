package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/dns-service/models"
)

// @Summary      intranet dns crud
// @Description  基于bind9 dns服务的针对dns的动态添加
// @Production   json
// @Production   json
// @Tags         dns-operation
// @Param        mock  query   string       false  "mock参数"
// @Success      200   object  StdResponse  "测试接口"
// @Router       /api/v1/dns [GET]
func addDns(c *gin.Context) {

}

func delDns(c *gin.Context) {

}

func updateDns(c *gin.Context) {

}

func getDns(c *gin.Context) {

}

func getDnsZone(c *gin.Context) {

}

func LoadDnsApis(r *gin.Engine) {
	apis := []models.Api{
		{Path: "/dns/record", Method: http.MethodGet, Description: "dns记录查询", Handler: getDns},
		{Path: "/dns", Method: http.MethodGet, Description: "dns查询", Handler: getDns},
		{Path: "/dns", Method: http.MethodDelete, Description: "dns删除", Handler: delDns},
		{Path: "/dns", Method: http.MethodPut, Description: "dns更新", Handler: updateDns},
		{Path: "/dns", Method: http.MethodPost, Description: "dns新增", Handler: addDns},

		{Path: "/dns/zone", Method: http.MethodGet, Description: "dns zone查询", Handler: getDnsZone},
	}
	loadApi(r, ginGroupApiV1, apis)
}
