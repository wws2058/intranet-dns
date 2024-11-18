package models

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tswcbyy1107/dns-service/database"
)

// api
type Api struct {
	BaseModel

	Path        string `gorm:"type:varchar(256) not null;index:idx_api_method" json:"path,omitempty"`  // api path
	Method      string `gorm:"type:varchar(16) not null;index:idx_api_method" json:"method,omitempty"` // api method
	Description string `gorm:"type:varchar(256) not null" json:"description,omitempty"`                // api description
	Active      bool   `gorm:"type:tinyint(1) default true" json:"active,omitempty"`                   // api is activated or not, 0 false
	Audit       bool   `gorm:"type:tinyint(1) default true" json:"audit,omitempty"`                    // api is audited or not, 0 false

	Handler gin.HandlerFunc `gorm:"-" json:"-"`
}

func (a *Api) TableName() string {
	return "apis"
}

// save gin api in db
func ApisRegisterInDB(apis []Api) {
	db := database.DB

	toBeSavedApis := []Api{}
	for _, api := range apis {
		tmpApi := &Api{}
		api.Path = strings.Replace(api.Path, "//", "/", -1)
		err := db.Model(&Api{}).Where("path = ? and method = ?", api.Path, api.Method).First(tmpApi).Error
		if err != nil {
			toBeSavedApis = append(toBeSavedApis, api)
		}
	}

	if len(toBeSavedApis) > 0 {
		db.Create(toBeSavedApis)
	}
}

// clean unused apis
func ApisCleanInDB(ginApis []Api) {
	db := database.DB

	dbApis := []Api{}
	toBeDeleteApisMap := make(map[string]uint)
	err := db.Model(&Api{}).Find(&dbApis).Error
	if err != nil {
		return
	}

	for _, api := range dbApis {
		toBeDeleteApisMap[fmt.Sprintf("%s/%s", api.Path, api.Method)] = api.Id
	}
	for _, api := range ginApis {
		delete(toBeDeleteApisMap, fmt.Sprintf("%s/%s", api.Path, api.Method))
	}

	ids := []uint{}
	if len(toBeDeleteApisMap) > 0 {
		for _, id := range toBeDeleteApisMap {
			ids = append(ids, id)
		}
	} else {
		return
	}

	TemplateSoftDelete(DaoDBReq{
		Dst: Api{},
		Where: map[string]interface{}{
			"id in (?)": ids,
		},
	})
}
