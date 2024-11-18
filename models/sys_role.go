package models

import (
	"github.com/tswcbyy1107/dns-service/database"
)

const (
	SuperAdmin = "super_admin" // all rw
	Admin      = "admin"       // partial high-risk operations
	CommonUser = "common_user" // common rw
	Guest      = "guest"       // partial r
)

type SysRole struct {
	BaseModel
	Name   string        `gorm:"type:varchar(32) not null;index:idx_name" json:"name,omitempty" binding:"len>0"`
	NameCn string        `gorm:"type:varchar(32) not null;index:idx_namecn" json:"name_cn,omitempty" binding:"len>0"`
	ApiIds mySlice[uint] `gorm:"type:varchar(1024)" json:"api_ids,omitempty" binding:"gt > 0"`

	AccessibleApis []Api `gorm:"-" json:"accessible_apis,omitempty"`
}

// role's accessible apis detail
func (r *SysRole) ApiDetails() (err error) {
	accessibleApis := []Api{}

	db := database.DB.Model(&Api{})
	if r.Name == SuperAdmin {
		db = db.Where("id > 0")
	} else {
		db = db.Where("id in (?)", r.ApiIds)
	}

	err = db.Find(&accessibleApis).Error
	if err != nil {
		return
	}
	r.AccessibleApis = accessibleApis
	return
}
