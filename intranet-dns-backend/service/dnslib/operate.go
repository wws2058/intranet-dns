package dnslib

import (
	"fmt"
	"strings"

	"github.com/wws2058/intranet-dns/database"
	"github.com/wws2058/intranet-dns/models"
	"github.com/wws2058/intranet-dns/utils"
)

func AddIntranetDns(rr *models.DnsRR, user string) (err error) {
	if err := rr.PreCheck(); err != nil {
		return err
	}

	data, err := IntranetRRQueryAll(rr.RecordName, rr.Zone)
	if err != nil {
		err = fmt.Errorf("dns query failed:%s", err.Error())
		return
	}
	existTypes := []string{}
	existRRs := []string{}
	for _, rr := range data {
		existTypes = append(existTypes, rr.RecordType)
		existRRs = append(existRRs, fmt.Sprintf("%v%v%v", rr.RecordName, rr.RecordType, rr.RecordContent))
	}

	newRRs, err := rr.ToRRs()
	if err != nil {
		return
	}
	for _, nr := range newRRs {
		nRecord, err := RRToDnsRR(nr, rr.Zone)
		if err != nil {
			return err
		}
		str := fmt.Sprintf("%v%v%v", nRecord.RecordName, nRecord.RecordType, nRecord.RecordContent)
		if utils.Contains(existRRs, str) {
			return fmt.Errorf("exist rr:%v", nr.String())
		}
	}

	existTypes = utils.RemoveRepeatedElement(existTypes)
	// check type: AAAA A conflict with CNAME
	if (rr.RecordType == models.AAAAType || rr.RecordType == models.AType) &&
		utils.Contains(existTypes, models.CnameType) {
		return fmt.Errorf("%s exist cname", rr.RecordName)
	}
	if rr.RecordType == models.CnameType &&
		(utils.Contains(existTypes, models.AType) || utils.Contains(existTypes, models.AAAAType)) {
		return fmt.Errorf("%s exist A/AAAA", rr.RecordName)
	}

	// rfc2136, some cases will ignore err, rcode is 0
	if err := IntranetDynamicDns(rr, nil, "add"); err != nil {
		return err
	}

	records := []*models.DnsRecord{}
	if rr.RecordType == "A" || rr.RecordType == "AAAA" {
		ips := strings.Split(rr.RecordContent, ",")
		for _, ip := range ips {
			dnsRR := models.DnsRR{
				RecordName:    rr.RecordName,
				RecordTtl:     rr.RecordTtl,
				Zone:          rr.Zone,
				RecordType:    rr.RecordType,
				RecordContent: ip,
			}

			records = append(records, &models.DnsRecord{
				DnsRR:   dnsRR,
				Creator: user,
			})
		}
	} else {
		records = append(records, &models.DnsRecord{
			DnsRR:   *rr,
			Creator: user,
		})
	}
	err = database.DB.CreateInBatches(records, len(records)).Error
	return
}

// clean=true remove all same type rrs; clean=false remove rr
func DelIntranetDns(id uint, clean bool) (err error) {
	record := &models.DnsRecord{}
	err = database.DB.Model(&models.DnsRecord{}).Where("id = ?", id).First(record).Error
	if err != nil {
		return
	}

	if !utils.Contains[string](models.LegalDnsType, record.RecordType) {
		err = fmt.Errorf("illegal record type:%s", record.RecordType)
		return
	}

	filter := map[string]interface{}{
		"id": id,
	}
	// Remove rr
	action := "del"
	if clean {
		// Remove rrs by name
		action = "del_all"
		filter = map[string]interface{}{
			"record_name": record.RecordName,
			"zone":        record.Zone,
			"record_type": record.RecordType,
		}
	}
	if err := IntranetDynamicDns(&record.DnsRR, nil, action); err != nil {
		return err
	}
	err = database.DB.Model(&models.DnsRecord{}).Where(filter).Update("deleted", 1).Error
	return
}

type UpdateDnsReq struct {
	Id            uint   `json:"id" binding:"required,gt=0"`
	RecordName    string `json:"record_name,omitempty"`
	RecordTtl     int    `json:"record_ttl,omitempty"`
	RecordContent string `json:"record_content,omitempty"`
}

func UpdateIntranetDns(req *UpdateDnsReq) (err error) {
	oldRecord := models.DnsRecord{}
	err = database.DB.Model(&models.DnsRecord{}).Where("id = ?", req.Id).First(&oldRecord).Error
	if err != nil {
		return
	}

	if !utils.Contains(models.LegalDnsType, oldRecord.RecordType) {
		err = fmt.Errorf("illegal record type:%s", oldRecord.RecordType)
		return
	}

	fields := []string{}
	newRecord := oldRecord
	if len(req.RecordName) > 0 {
		newRecord.RecordName = req.RecordName
		fields = append(fields, "record_name")
	}
	if req.RecordTtl > 0 {
		newRecord.RecordTtl = req.RecordTtl
		fields = append(fields, "record_ttl")
	}
	if len(req.RecordContent) > 0 {
		if oldRecord.RecordType == models.CnameType {
			newRecord.RecordContent = req.RecordContent
		}
		if len(strings.Split(req.RecordContent, ",")) > 1 {
			return fmt.Errorf("batch updates are not allowed")
		}
		newRecord.RecordContent = req.RecordContent
		fields = append(fields, "record_content")
	}
	if len(fields) == 0 {
		return nil
	}
	if err := newRecord.PreCheck(); err != nil {
		return err
	}
	err = IntranetDynamicDns(&oldRecord.DnsRR, &newRecord.DnsRR, "modify")
	if err != nil {
		return
	}
	err = models.TemplateUpdate(&newRecord, fields)
	return
}
