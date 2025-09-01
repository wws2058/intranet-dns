package models

import (
	"github.com/wws2058/intranet-dns/database"
	"gorm.io/gorm"
)

const (
	SuperAdmin = "super_admin" // all rw
	Admin      = "admin"       // partial high-risk operations
	CommonUser = "common_user" // common rw
	Guest      = "guest"       // partial r
)

type SysRole struct {
	BaseModel
	Name   string        `gorm:"type:varchar(32) not null;uniqueIndex:uk_name" json:"name,omitempty" binding:"gt=0"`
	NameCn string        `gorm:"type:varchar(32) not null;index:idx_namecn" json:"name_cn,omitempty" binding:"gt=0"`
	ApiIds MySlice[uint] `gorm:"type:varchar(1024)" json:"api_ids,omitempty" binding:"gt=0"`

	AccessibleApis []Api `gorm:"-" json:"accessible_apis,omitempty"`
}

// role's accessible apis detail
func (r *SysRole) ApiDetails() (err error) {
	accessibleApis := []Api{}

	db := database.DB.Model(&Api{})
	if r.Name == SuperAdmin {
		db = db.Where("id > 0")
	} else {
		db = db.Where("id in (?)", []uint(r.ApiIds))
	}

	err = db.Find(&accessibleApis).Error
	if err != nil {
		return
	}
	r.AccessibleApis = accessibleApis
	return
}

func DelSysRole(roleID uint) error {
	db := database.DB

	err := db.Transaction(func(tx *gorm.DB) error {
		// soft del role
		err := db.Model(&SysRole{}).Where("id = ?", roleID).Update("deleted", true).Error
		if err != nil {
			return err
		}

		// update user
		users := []SysUser{}
		err = db.Where(ColumnContains("role_ids"), roleID).Find(&users).Error
		if err != nil {
			return err
		}
		for _, u := range users {
			err := database.DB.Model(&u).Update("role_ids", u.RoleIds.Del(roleID)).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
