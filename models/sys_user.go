package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tswcbyy1107/dns-service/config"
	"github.com/tswcbyy1107/dns-service/database"
	"github.com/tswcbyy1107/dns-service/lib/redis"
)

type SysUser struct {
	BaseModel
	Name       string        `gorm:"type:varchar(32) not null; uniqueIndex:uk_name" json:"name,omitempty" binding:"min=1"` // user name
	NameCn     string        `gorm:"type:varchar(32) not null; index:idx_namecn" json:"name_cn,omitempty" binding:"min=1"` // user cn name
	Email      string        `gorm:"type:varchar(64) not null" json:"email,omitempty" binding:"min=1"`                     // user email address
	Active     bool          `gorm:"type:tinyint(1) default true" json:"active,omitempty"`                                 // user is banned, active=0
	LastLogin  JsonTime      `gorm:"type:datetime" json:"last_login,omitempty"`                                            // user last login at
	LoginTimes int           `gorm:"type:int(10) default 0" json:"login_times,omitempty"`                                  // user login times
	RoleIds    mySlice[uint] `gorm:"type:varchar(1024)" json:"role_ids,omitempty"`                                         // user roles
}

// get role accessible active api path and method, subitem = path/method
func (u *SysUser) GetAccessibleApis() ([]string, error) {
	pathWithMethod := []string{}

	// try to get cache
	key := fmt.Sprintf("%s_%s_apis", config.GlobalConfig.App.Name, u.Name)
	if err := redis.LoadCache(key, &pathWithMethod); err == nil {
		return pathWithMethod, nil
	}

	roles := []SysRole{}
	err := database.DB.Model(&SysRole{}).Where("id in (?)", u.RoleIds).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	apiIDs := []uint{}
	for _, r := range roles {
		apiIDs = append(apiIDs, r.ApiIds...)
	}

	apis := []Api{}
	err = database.DB.Model(&Api{}).Where("id in (?)", apiIDs).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	// ignore inactivated apis
	for _, api := range apis {
		if !api.Active {
			continue
		}
		pathWithMethod = append(pathWithMethod, fmt.Sprintf("%s/%s", api.Path, api.Method))
	}

	// cache user's apis
	bytes, _ := json.Marshal(pathWithMethod)
	redis.Cache(key, string(bytes), time.Minute*5)

	return pathWithMethod, nil
}
