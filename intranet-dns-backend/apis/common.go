package apis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wws2058/intranet-dns/models"
)

// gin group
var ginGroupApiV1 = "/api/v1"

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
