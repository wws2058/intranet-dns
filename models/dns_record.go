package models

type DnsRecord struct {
	BaseModel

	RecordName    string `gorm:"type:varchar(256) not null;index:idx_record_name" json:"record_name,omitempty"`       // domain
	RecordTtl     int    `gorm:"type:int(10) not null default 60" json:"record_ttl,omitempty"`                        // ttl
	Zone          string `gorm:"type:varchar(256) not null;index:idx_zone" json:"zone,omitempty"`                     // zone
	RecordType    string `gorm:"type:varchar(16) not null" json:"record_type,omitempty"`                              // type
	RecordContent string `gorm:"type:varchar(256) not null;index:idx_record_content" json:"record_content,omitempty"` // content

	Service string `gorm:"type:varchar(256);index:idx_service" json:"services,omitempty"` // domain's service
	Creator string `gorm:"type:varchar(32)" json:"creator,omitempty"`                     // domain's user
}

// bind dns rr in db
func (d *DnsRecord) TableName() string {
	return "dns_records"
}
