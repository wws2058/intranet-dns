package models

// api request audit log
type AuditLog struct {
	BaseModel
	UserName     string `gorm:"type:varchar(64) not null; index:idx_user_name" json:"user_name,omitempty" form:"user_name"` // user name(web) or appid(api)
	RequestID    string `gorm:"type:varchar(64) not null; index:idx_request_id" json:"request_id,omitempty" form:"request_id"`
	ClientIP     string `gorm:"type:varchar(64) not null; index:idx_client_ip" json:"client_ip,omitempty" form:"client_ip"`
	URL          string `gorm:"type:varchar(1024) not null" json:"url,omitempty"`
	Method       string `gorm:"type:varchar(16) not null" json:"method,omitempty"` // api method
	RequestBody  string `gorm:"type:text not null" json:"request_body,omitempty"`
	ResponseBody string `gorm:"type:text not null" json:"response_body,omitempty"`
	TimeCost     int    `gorm:"type:int(10) not null" json:"time_cost,omitempty"`
}
