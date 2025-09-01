package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wws2058/intranet-dns/config"
	"github.com/wws2058/intranet-dns/database"
	"github.com/wws2058/intranet-dns/service/redis"
	"github.com/wws2058/intranet-dns/utils"
	"gorm.io/gorm"
)

type SysUser struct {
	BaseModel
	Name       string        `gorm:"type:varchar(32) not null; uniqueIndex:uk_name" json:"name,omitempty" binding:"min=1"` // user name
	Password   string        `gorm:"type:varchar(128) not null; index: idx_password" json:"-"`                             // user password, sha256 encode save
	NameCn     string        `gorm:"type:varchar(32) not null; index:idx_namecn" json:"name_cn,omitempty" binding:"min=1"` // user cn name
	Email      string        `gorm:"type:varchar(64) not null" json:"email,omitempty" binding:"min=1"`                     // user email address
	Active     bool          `gorm:"type:tinyint(1) default true" json:"active,omitempty"`                                 // user is banned, active=0
	LastLogin  JsonTime      `gorm:"type:datetime" json:"last_login,omitempty"`                                            // user last login at
	LoginTimes int           `gorm:"type:int(10) default 0" json:"login_times,omitempty"`                                  // user login times
	RoleIds    MySlice[uint] `gorm:"type:varchar(1024)" json:"role_ids,omitempty"`                                         // user roles
}

func (u *SysUser) IsSuperAdmin() bool {
	superAdmin := &SysRole{}
	err := database.DB.Where("name = ?", SuperAdmin).First(superAdmin).Error
	if err != nil {
		return false
	}
	if u.RoleIds.Contains(superAdmin.Id) {
		return true
	}
	return false
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
	redis.Cache(key, string(bytes), time.Minute*1)

	return pathWithMethod, nil
}

// sha256 user's password before create user
func (u *SysUser) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.Password) == 0 || utils.IsSha256(u.Password) {
		return
	}
	u.Password = utils.Sha256Hash(utils.UserPasswdSalt + u.Password)
	return
}

// sha256 user's password before update user's password
func (u *SysUser) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(u.Password) == 0 || utils.IsSha256(u.Password) {
		return
	}
	u.Password = utils.Sha256Hash(utils.UserPasswdSalt + u.Password)
	return
}

// sha256 user's password
func (u *SysUser) Sha256Password() {
	if len(u.Password) == 0 || utils.IsSha256(u.Password) {
		return
	}
	u.Password = utils.Sha256Hash(utils.UserPasswdSalt + u.Password)
}

func UpdateUserLoginInfo(username string) {
	database.DB.Model(&SysUser{}).Where("name = ?", username).Updates(map[string]interface{}{
		"login_times": gorm.Expr("login_times + 1"),
		"last_login":  JsonTime(time.Now()),
	})
}
