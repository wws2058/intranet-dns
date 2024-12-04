package models

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/tswcbyy1107/intranet-dns/database"
	"gorm.io/plugin/soft_delete"
)

type PageReq struct {
	Page     int `json:"page" form:"page" binding:"gte=1"`                     // page, start at 1
	PageSize int `json:"page_size" form:"page_size" binding:"gte=10,lte=1000"` // page size, 10<=x<=1000
}

// page query offset
func (r *PageReq) Offset() int {
	if r == nil {
		r = new(PageReq)
		r.Page = 1
		r.PageSize = 20
	}

	if (r.Page-1)*r.PageSize <= 0 {
		return 0
	}
	return (r.Page - 1) * r.PageSize
}

type PageRsp struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// base table model
type BaseModel struct {
	Id        uint                  `gorm:"primaryKey;autoIncrement" json:"id"`                                  // primary id
	CreatedAt JsonTime              `gorm:"type:datetime;autoCreateTime;index:idx_created_at" json:"created_at"` // create time
	UpdatedAt JsonTime              `gorm:"type:datetime;autoUpdateTime;index:idx_updated_at" json:"updated_at"` // modified time
	Deleted   soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-" `                                           // record is deleted, 0 false
}

// common db dao request
type DaoDBReq struct {
	PageReq PageReq // page req
	PageRsp PageRsp // page rsp

	Dst interface{} // db models

	ModelFilter interface{}            // conditions in =
	Where       map[string]interface{} // custom conditions, eg: 'create_at > ?':'xxxx'

	OrderBy string // order condition, eg: "create_time desc, update_time"
}

// common db query
func TemplateQuery(query *DaoDBReq) (err error) {
	model, isSlice := newModel(query.Dst)
	db := database.DB.Model(model)
	if query.ModelFilter != nil {
		db = db.Where(query.ModelFilter)
	}
	for k, v := range query.Where {
		db = db.Where(k, v)
	}
	if isSlice {
		err = db.Find(query.Dst).Error
	} else {
		err = db.First(query.Dst).Error
	}
	return
}

// common db page query
func TemPlatePageQuery(query *DaoDBReq) (err error) {
	var total int64

	model, _ := newModel(query.Dst)
	db := database.DB.Model(model)

	if query.ModelFilter != nil {
		db = db.Where(query.ModelFilter)
	}
	for k, v := range query.Where {
		db = db.Where(k, v)
	}
	if query.OrderBy != "" {
		db = db.Order(query.OrderBy)
	} else {
		db = db.Order("id desc")
	}
	db.Count(&total)
	if err := db.Offset(query.PageReq.Offset()).Limit(query.PageReq.PageSize).Find(query.Dst).Error; err != nil {
		return err
	}

	query.PageRsp = PageRsp{
		Page:     query.PageReq.Page,
		PageSize: query.PageReq.PageSize,
		Total:    int(total),
	}
	return nil
}

// common save: create or update
func TemplateCreate(opt interface{}) (err error) {
	return database.DB.Save(opt).Error
}

// common update: field + conditions to update zero fields, primary key required
func TemplateUpdate(model interface{}, fieldsToUpdate []string) (err error) {
	if len(fieldsToUpdate) == 0 {
		return nil
	}
	return database.DB.Model(model).Select(fieldsToUpdate).Updates(model).Error
}

// common delete single record: deleted column and primary key are required,
func TemplateSoftDelete(opt DaoDBReq) (err error) {
	db := database.DB.Model(opt.Dst)
	if opt.ModelFilter != nil {
		db = db.Where(opt.ModelFilter)
	}
	for k, v := range opt.Where {
		db = db.Where(k, v)
	}
	return db.Update("deleted", 1).Error
}

// get specific model
func newModel(src interface{}) (interface{}, bool) {
	object := reflect.TypeOf(src)
	isSlice := false

	// \*T -> T  \*[] -> []
	if object.Kind() == reflect.Ptr {
		object = object.Elem()
	}

	// [item] -> item
	if object.Kind() == reflect.Slice {
		object = object.Elem()
		isSlice = true
	}

	// \*T -> T
	if object.Kind() == reflect.Ptr {
		object = object.Elem()
	}
	return reflect.New(object).Interface(), isSlice
}

// auto create tables
func AutoMigrate() {
	err := database.DB.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		&Api{},
		&AuditLog{},
		&SysRole{},
		&SysUser{},
		&Cronjob{},
		&DnsRecord{},
		&DnsZone{},
	)

	// create demo role
	// superAdminR := &SysRole{
	// 	Name:   SuperAdmin,
	// 	NameCn: "超级管理员",
	// }
	// database.DB.Create(superAdminR)
	// create demo user
	// superUsers := &SysUser{
	// 	Name:     "somebody",
	// 	NameCn:   "系统管理员",
	// 	Email:    "china.qq.com",
	// 	Password: "12345678",
	// 	RoleIds:  MySlice[uint]{superAdminR.Id},
	// 	Active:   true,
	// }
	// database.DB.Create(superUsers)
	// create demo dns zone
	// testZone := &DnsZone{
	// 	Zone:        "test.com.",
	// 	NsAddress:   "1.1.1.1:53",
	// 	TsigName:    "key-name.",
	// 	TsigSecret:  "secret",
	// 	Description: "test zone",
	// 	Creator:     "somebody",
	// }
	// database.DB.Create(testZone)

	if err != nil {
		logrus.WithField("mysql", "auto_migrate").Error(err)
	} else {
		logrus.WithField("mysql", "auto_migrate").Info("succeed")
	}
}

func ColumnContains(column string) string {
	return fmt.Sprintf("find_in_set(?, %v)", column)
}
