package models

import "github.com/miekg/dns"

// note: https://www.cnblogs.com/RichardLuo/p/DNS_P3.html

// dns zone info, FQDN format
type DnsZone struct {
	BaseModel

	Zone        string `gorm:"type:varchar(64) not null;uniqueIndex:uk_ns_zone" json:"zone,omitempty"`       // ns name zone name FQDN
	NsAddress   string `gorm:"type:varchar(64) not null;uniqueIndex:uk_ns_zone" json:"ns_address,omitempty"` // ns server ip:port
	TsigName    string `gorm:"type:varchar(64) not null;index:idx_tsig_name" json:"tsig_name,omitempty"`     // tsig key name
	TsigSecret  string `gorm:"type:varchar(64) not null" json:"-"`                                           // ns dynamic update key, to be encoded
	Description string `gorm:"type:varchar(128) default null" json:"description,omitempty"`                  // zone description
	Creator     string `gorm:"type:varchar(64) not null;index:idx_creator" json:"creator,omitempty"`         // zone creator
}

func (d *DnsZone) SetFqdn() {
	d.TsigName = dns.Fqdn(d.TsigName)
	d.Zone = dns.Fqdn(d.Zone)
}
