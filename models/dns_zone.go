package models

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// note: https://www.cnblogs.com/RichardLuo/p/DNS_P3.html

// dns zone info, FQDN format, intranet second-level zone, eg test.com. company.net.  funny.cn.
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

func (d *DnsZone) PreCheck() (err error) {
	d.SetFqdn()
	err = CheckZone(d.Zone)
	return
}

// intranet second-level zone, eg test.com. company.net.  funny.cn.
// zone + address uk column in db
func CheckZone(zone string) (err error) {
	zone = dns.Fqdn(zone)
	if len(strings.Split(zone, ".")) != 3 {
		return fmt.Errorf("illegal second-level zone formal")
	}
	return
}
