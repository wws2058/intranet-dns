package dnslib

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
	"github.com/tswcbyy1107/intranet-dns/models"
)

/*
note: https://datatracker.ietf.org/doc/html/rfc2136#autoid-2
RR: single dns record, test.com in A 1.1.1.1
RRset: a  collection of RR, test.com in A 1.1.1.1 + test.com in A 2.2.2.2

github: https://pkg.go.dev/github.com/miekg/dns@v1.1.48#section-readme
3.2.4 - Table Of Metavalues Used In Prerequisite Section

 CLASS    TYPE     RDATA    Meaning                    Function
 --------------------------------------------------------------
 ANY      ANY      empty    Name is in use             dns.NameUsed
 ANY      rrset    empty    RRset exists (value indep) dns.RRsetUsed
 NONE     ANY      empty    Name is not in use         dns.NameNotUsed
 NONE     rrset    empty    RRset does not exist       dns.RRsetNotUsed
 zone     rrset    rr       RRset exists (value dep)   dns.Used

 3.4.2.6 - Table Of Metavalues Used In Update Section

 CLASS    TYPE     RDATA    Meaning                     Function
 ---------------------------------------------------------------
 ANY      ANY      empty    Delete all RRsets from name dns.RemoveName
 ANY      rrset    empty    Delete an RRset             dns.RemoveRRset
 NONE     rrset    rr       Delete an RR from RRset     dns.Remove
 zone     rrset    rr       Add to an RRset             dns.Insert
*/

// intranet dns crud with bind dns server, action: add del modify del_all, timeout in 3s
// CNAME, A , AAAA; AAAA A conflict CNAME
func IntranetDynamicDns(record *models.DnsRR, updateRecord *models.DnsRR, action string) (err error) {
	dnsZone := &models.DnsZone{}
	if err := models.TemplateQuery(&models.DaoDBReq{
		Dst:         dnsZone,
		ModelFilter: map[string]string{"zone": record.Zone},
	}); err != nil {
		return err
	}
	dnsZone.SetFqdn()

	rrs, err := record.ToRRs()
	if err != nil {
		return
	}
	msg := &dns.Msg{}
	msg.SetUpdate(dnsZone.Zone)
	msg.SetTsig(dnsZone.TsigName, dns.HmacSHA256, 300, time.Now().Unix())
	// msg.RRsetUsed(rrs)
	switch action {
	case "add":
		// if record.RecordType == models.CnameType {
		// 	msg.NameNotUsed(rrs)
		// }
		// if record.RecordType == models.AAAAType || record.RecordType == models.AType {
		// 	msg.RRsetNotUsed(rrs)
		// }
		msg.Insert(rrs)
	case "del":
		msg.Remove(rrs)
	case "del_all":
		msg.RemoveName(rrs)
	case "modify":
		upRRs, err := updateRecord.ToRRs()
		if err != nil {
			return err
		}
		msg.Remove(rrs)
		msg.Insert(upRRs)
	default:
		return fmt.Errorf("unsupported action:%s", action)
	}

	dnsClient := &dns.Client{
		TsigSecret: map[string]string{
			dnsZone.TsigName: dnsZone.TsigSecret,
		},
		Timeout: 3 * time.Second,
	}

	dnsRsp, rtt, err := dnsClient.Exchange(msg, dnsZone.NsAddress)
	if dnsRsp != nil && dnsRsp.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("dns rcode err:%v", dnsRsp.Rcode)
	}
	logrus.WithField("dynamic dns", action).Infof("%v %v %v %v", *record, rtt, dnsRsp.Rcode, err)
	return
}
