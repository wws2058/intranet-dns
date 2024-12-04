package models

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/miekg/dns"
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

// db dns record to dns dynamic dns rr, A AAAA CNMAE, A AAAA multiple content separated by ','
func (d *DnsRR) ToDnsRRs() (rrs []dns.RR, err error) {
	d.SetFqdn()

	rrHeader := dns.RR_Header{
		Name:  d.RecordName,
		Class: dns.ClassINET,
		Ttl:   uint32(d.RecordTtl),
	}
	rrs = []dns.RR{}
	switch d.RecordType {
	case "A":
		rrHeader.Rrtype = dns.TypeA
		ips := strings.Split(d.RecordContent, ",")
		for _, ipv4 := range ips {
			rr := &dns.A{
				A:   net.ParseIP(ipv4),
				Hdr: rrHeader,
			}
			rrs = append(rrs, rr)
		}
	case "CNAME":
		rrHeader.Rrtype = dns.TypeCNAME
		rrs = append(rrs, &dns.CNAME{
			Target: d.RecordContent,
			Hdr:    rrHeader,
		})
	case "AAAA":
		rrHeader.Rrtype = dns.TypeAAAA
		ips := strings.Split(d.RecordContent, ",")
		for _, ipv6 := range ips {
			rr := &dns.AAAA{
				AAAA: net.ParseIP(ipv6),
				Hdr:  rrHeader,
			}
			rrs = append(rrs, rr)
		}
	default:
		err = fmt.Errorf("unsupported type:%s", d.RecordType)
	}
	return rrs, err
}

// intranet dns format: {name}.{zone}, name cannot contain dots [a-zA-Z0-9._+-]+$, check before add or update
func (d *DnsRR) PreCheck() bool {
	d.SetFqdn()
	re := regexp.MustCompile(`^[a-zA-Z0-9._+-]+$`)
	if !re.MatchString(d.RecordName) {
		return false
	}
	if !strings.HasSuffix(d.RecordName, d.Zone) {
		return false
	}
	rName := strings.TrimRight(d.RecordName, d.Zone)
	return len(strings.Split(rName, ".")) == 1
}
