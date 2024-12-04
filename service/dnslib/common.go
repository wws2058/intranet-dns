package dnslib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/tswcbyy1107/intranet-dns/models"
)

type EdnsRRs struct {
	ISP    string         `json:"isp,omitempty"`
	DnsRRs []models.DnsRR `json:"dns_rrs,omitempty"`
}

type ProvinceIspDns struct {
	ISP   string `json:"isp,omitempty"`    // metadata {province}-{isp}
	DnsIP string `json:"dns_ip,omitempty"` // ns {ip}
}

// dns ip address of isps in each province captial in China
var PublicDnsIP = []ProvinceIspDns{
	{
		ISP:   "北京-电信",
		DnsIP: "219.141.136.10",
	},
	{
		ISP:   "北京-移动",
		DnsIP: "221.130.33.60",
	},
	{
		ISP:   "北京-联通",
		DnsIP: "202.106.0.20",
	},
	{
		ISP:   "上海-电信",
		DnsIP: "202.96.209.5",
	},
	{
		ISP:   "上海-移动",
		DnsIP: "211.136.150.66",
	},
	{
		ISP:   "上海-联通",
		DnsIP: "210.22.70.3",
	},
	{
		ISP:   "深圳-电信",
		DnsIP: "202.96.134.133",
	},
	{
		ISP:   "深圳-移动",
		DnsIP: "211.136.192.6",
	},
	{
		ISP:   "深圳-联通",
		DnsIP: "210.21.196.6",
	},
}

// A -> 1, AAAA -> 28, CNAME -> 5
func RTypeStrToUint(strType string) (rrType uint16) {
	switch strType {
	case "A":
		rrType = dns.TypeA
	case "CNAME":
		rrType = dns.TypeCNAME
	case "PTR":
		rrType = dns.TypePTR
	case "MX":
		rrType = dns.TypeMX
	case "TXT":
		rrType = dns.TypeTXT
	case "SRV":
		rrType = dns.TypeSRV
	case "AAAA":
		rrType = dns.TypeAAAA
	case "NS":
		rrType = dns.TypeNS
	}
	return rrType
}

// dns rr to dns record in db
func RRToDnsRR(rr dns.RR, zone string) (record models.DnsRR, err error) {
	rrArray := strings.Split(rr.String(), "\t")
	if len(rrArray) != 5 {
		err = fmt.Errorf("rr format err: %s", rr.String())
		return
	}
	ttl, _ := strconv.Atoi(rrArray[1])
	record = models.DnsRR{
		RecordName:    rrArray[0],
		RecordTtl:     ttl,
		RecordType:    rrArray[3],
		RecordContent: rrArray[4],
		Zone:          dns.Fqdn(zone),
	}
	return
}
