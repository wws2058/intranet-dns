package dnslib

import (
	"fmt"

	"github.com/tswcbyy1107/intranet-dns/database"
	"github.com/tswcbyy1107/intranet-dns/models"
	"gorm.io/gorm"
)

func AddDns(record *models.DnsRR) (err error) {
	if !record.PreCheck() {
		return fmt.Errorf("illegal intranet dns format")
	}

	// query
	rrs, err := IntranetRRQuery(record.RecordName, record.Zone, record.RecordType)
	if err != nil {
		return
	}
	if len(rrs) != 0 {
		return fmt.Errorf("records exists")
	}

	db := database.DB

	db.Transaction(func(tx *gorm.DB) (err error) {

		return
	})

	return
}
