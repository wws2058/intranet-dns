package models

import (
	"net"
	"sort"

	"github.com/miekg/dns"
)

const (
	ProbeSucceed     = 0
	ProbeQueryFailed = 1
	ProbeMisMatch    = 2
)

type DnsProbe struct {
	BaseModel

	RecordName   string          `gorm:"type:varchar(256) not null;index:idx_record_name" json:"record_name,omitempty"` // domain FQDN
	Zone         string          `gorm:"type:varchar(256) not null;index:idx_zone" json:"zone,omitempty"`               // zone FQDN
	ExpectAnswer MySlice[string] `gorm:"type:varchar(1024) not null" json:"expect_answer,omitempty"`                    // expected rrs, ips
	Creator      string          `gorm:"type:varchar(64) not null;index:idx_creator" json:"creator,omitempty"`          // probe creator
	Result       uint            `gorm:"type:smallint unsigned default 0" json:"succeed"`                               // last probe result
	Intranet     bool            `gorm:"type:tinyint(1) default true" json:"intranet"`                                  // intranet dns
	// RecordType   string `gorm:"type:varchar(16) not null" json:"record_type,omitempty"`                        // type
}

func (d *DnsProbe) SetFqdn() {
	d.RecordName = dns.Fqdn(d.RecordName)
	d.Zone = dns.Fqdn(d.Zone)
	fqdn := MySlice[string]{}
	for _, v := range d.ExpectAnswer {
		ip := net.ParseIP(v)
		if ip == nil {
			v = dns.Fqdn(v)
		}
		fqdn = append(fqdn, v)
	}
	sort.Sort(fqdn)
	d.ExpectAnswer = fqdn
}
