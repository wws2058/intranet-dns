package models

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
	"github.com/tswcbyy1107/intranet-dns/utils"
)

/*
note:
https://datatracker.ietf.org/doc/html/rfc1035
https://datatracker.ietf.org/doc/html/rfc2136
*/

// dns record type
const (
	AType     = "A"
	AAAAType  = "AAAA"
	CnameType = "CNAME"
	NsType    = "NS"
	SoaType   = "SOA"
	MxType    = "MX"
)

var (
	LegalDnsType = []string{AType, AAAAType, CnameType}
)

// dns rr in db , FQDN format
type DnsRecord struct {
	BaseModel

	DnsRR

	Creator   string `gorm:"type:varchar(64) not null;index:idx_creator" json:"creator,omitempty"` // domain's user
	ExtraInfo string `gorm:"type:varchar(1024) default null" json:"extra_info,omitempty"`          // domain's extra info
}

// dns rr , FQDN format
type DnsRR struct {
	RecordName    string `gorm:"type:varchar(256) not null;index:idx_record_name" json:"record_name,omitempty"`       // domain FQDN
	RecordTtl     int    `gorm:"type:int(10) not null default 60" json:"record_ttl,omitempty"`                        // ttl
	Zone          string `gorm:"type:varchar(256) not null;index:idx_zone" json:"zone,omitempty"`                     // zone FQDN
	RecordType    string `gorm:"type:varchar(16) not null" json:"record_type,omitempty"`                              // type
	RecordContent string `gorm:"type:varchar(256) not null;index:idx_record_content" json:"record_content,omitempty"` // content
}

func (d *DnsRR) SetFqdn() {
	d.RecordName = dns.Fqdn(d.RecordName)
	d.Zone = dns.Fqdn(d.Zone)
	if d.RecordType == CnameType {
		d.RecordContent = dns.Fqdn(d.RecordContent)
	}
	if d.RecordTtl == 0 {
		d.RecordTtl = 60
	}
}

// db dns record to dns dynamic dns rr with check, A AAAA CNMAE, A AAAA multiple content separated by ','
func (d *DnsRR) ToRRs() (rrs []dns.RR, err error) {
	err = d.PreCheck()
	if err != nil {
		return
	}

	if d.RecordType == AAAAType || d.RecordType == AType {
		ips := utils.RemoveRepeatedElement(strings.Split(d.RecordContent, ","))
		for _, ip := range ips {
			recordStr := fmt.Sprintf("%v %v IN %v %v", d.RecordName, d.RecordTtl, d.RecordType, ip)
			rr, err := dns.NewRR(recordStr)
			if err != nil {
				return nil, err
			}
			rrs = append(rrs, rr)
		}
	} else {
		recordStr := fmt.Sprintf("%v %v IN %v %v", d.RecordName, d.RecordTtl, d.RecordType, d.RecordContent)
		rr, err := dns.NewRR(recordStr)
		if err != nil {
			return nil, err
		}
		rrs = append(rrs, rr)
	}
	return rrs, err
}

// intranet dns format: {name}.{zone}, name cannot contain dots [a-zA-Z0-9._+-]+$, check before add or update
// set Fqdn
func (d *DnsRR) PreCheck() (err error) {
	if d.RecordType == CnameType {
		ip := net.ParseIP(d.RecordContent)
		if ip != nil {
			return fmt.Errorf("ip in CNAME type %v", ip.String())
		}
	}

	d.SetFqdn()

	if !utils.Contains(LegalDnsType, d.RecordType) {
		return fmt.Errorf("illegal record type")
	}
	if err := CheckZone(d.Zone); err != nil {
		return err
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._+-]+$`)
	if !re.MatchString(d.RecordName) {
		err = fmt.Errorf("illegal record name str")
		return
	}
	if !strings.HasSuffix(d.RecordName, d.Zone) {
		err = fmt.Errorf("record name and zonedo not match")
		return
	}
	return
}
