package dnslib

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/tswcbyy1107/intranet-dns/models"
	"golang.org/x/sync/errgroup"
)

// query intranet dns A AAAA CNAME rr
func IntranetRRQueryAll(domain, zone string) (rrs []models.DnsRR, err error) {
	domain = dns.Fqdn(domain)
	zone = dns.Fqdn(zone)

	g := new(errgroup.Group)
	ch := make(chan []models.DnsRR, len(models.LegalDnsType))
	for _, t := range models.LegalDnsType {
		g.Go(func() (err error) {
			rrs, err := IntranetRRQuery(domain, zone, t)
			if err != nil {
				return
			}
			ch <- rrs
			return
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(ch)
	for data := range ch {
		rrs = append(rrs, data...)
	}
	return
}

// intranet dns query, xxx.baidu.com baidu.com A, without tsig
func IntranetRRQuery(domain, zone, rtype string) (rrs []models.DnsRR, err error) {
	domain = dns.Fqdn(domain)
	zone = dns.Fqdn(zone)
	rrType := RTypeStrToUint(rtype)

	// get zone
	dnsZone := &models.DnsZone{}
	if err := models.TemplateQuery(&models.DaoDBReq{
		Dst:         dnsZone,
		ModelFilter: map[string]string{"zone": zone},
	}); err != nil {
		return nil, err
	}

	msg := &dns.Msg{}
	// msg.SetTsig(dnsZone.TsigName, dns.HmacSHA256, 300, time.Now().Unix())
	dnsClient := &dns.Client{
		Timeout: 3 * time.Second,
		// TsigSecret: map[string]string{
		// 	dnsZone.TsigName: dnsZone.TsigSecret,
		// },
	}
	msg.SetQuestion(domain, rrType)
	dnsRsp, _, err := dnsClient.Exchange(msg, dnsZone.NsAddress)
	if err != nil {
		return
	}
	if dnsRsp == nil {
		err = fmt.Errorf("no answer")
		return
	}
	// if dnsRsp.Rcode != dns.RcodeSuccess {
	// 	rcode = dnsRsp.Rcode
	// 	err = fmt.Errorf("dns failed rcode:%v", dns.RcodeToString[dnsRsp.Rcode])
	// 	return
	// }
	for _, rr := range dnsRsp.Answer {
		if rr.Header().Rrtype != rrType {
			continue
		}
		dnsRR, err := RRToDnsRR(rr, zone)
		if err != nil {
			continue
		}
		rrs = append(rrs, dnsRR)
	}
	return
}

// intranet axfr transfer in zone, get all records
func IntranetRRsInZone(zone string) (rrs []models.DnsRR, err error) {
	zone = dns.Fqdn(zone)

	dnsZone := &models.DnsZone{}
	if err := models.TemplateQuery(&models.DaoDBReq{
		Dst:         dnsZone,
		ModelFilter: map[string]string{"zone": zone},
	}); err != nil {
		return nil, err
	}

	msg := &dns.Msg{}
	msg.SetTsig(dnsZone.TsigName, dns.HmacSHA256, 300, time.Now().Unix())
	msg.SetAxfr(zone)
	t := &dns.Transfer{
		TsigSecret: map[string]string{
			dnsZone.TsigName: dnsZone.TsigSecret,
		},
		ReadTimeout: 1 * time.Second,
	}

	dnsRsp, err := t.In(msg, dnsZone.NsAddress)
	if err != nil {
		return
	}
	for c := range dnsRsp {
		for _, rr := range c.RR {
			dnsRR, err := RRToDnsRR(rr, zone)
			if err != nil {
				continue
			}
			rrs = append(rrs, dnsRR)
		}
	}
	return
}

// public dns 119.29.29.29 edns ipv4 query, give the remote nameserver an idea of where the client lives.
func PublicEDnsQueryRR(domain, clientIp string) (rrs []models.DnsRR, err error) {
	domain = dns.Fqdn(domain)
	// public dns query rrType does not take effect, require A return all
	rrType := RTypeStrToUint("A")

	opt := &dns.OPT{
		Hdr: dns.RR_Header{
			Name:   ".",
			Rrtype: dns.TypeOPT,
		},
		Option: []dns.EDNS0{
			&dns.EDNS0_SUBNET{
				Code:          dns.EDNS0SUBNET,
				Address:       net.ParseIP(clientIp).To4(),
				Family:        1,
				SourceNetmask: 32,
			},
		},
	}
	msg := &dns.Msg{}
	// public dns query rrType does not take effect
	msg.SetQuestion(domain, rrType)
	msg.Extra = append(msg.Extra, opt)

	dnsClient := &dns.Client{
		Timeout: 2 * time.Second,
	}

	dnsRsp, _, err := dnsClient.Exchange(msg, "119.29.29.29:53")
	if err != nil {
		return
	}
	if dnsRsp == nil {
		err = fmt.Errorf("no answer")
		return
	}
	if dnsRsp.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("dns failed rcode:%v", dns.RcodeToString[dnsRsp.Rcode])
		return
	}
	for _, rr := range dnsRsp.Answer {
		// if rr.Header().Rrtype != rrType {
		// 	continue
		// }
		dnsRR, err := RRToDnsRR(rr, ".")
		if err != nil {
			continue
		}
		rrs = append(rrs, dnsRR)
	}
	return
}
