package cronjob

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tswcbyy1107/intranet-dns/config"
	"github.com/tswcbyy1107/intranet-dns/database"
	"github.com/tswcbyy1107/intranet-dns/models"
	"github.com/tswcbyy1107/intranet-dns/service/dnslib"
)

var internalFunctionMaps = map[string]func() error{
	"test_function": testFunction,
	"sync_rrs":      SyncDBRecordWithNsRR,
	"dns_probe":     dnsProbe,
}

// get internal cronjob function names
func GetInternalFunctions() []string {
	names := []string{}
	for k := range internalFunctionMaps {
		names = append(names, k)
	}
	return names
}

// test cronjob function
func testFunction() error {
	rand.NewSource(time.Now().UnixNano())
	randomInt := rand.Intn(100)
	if randomInt%2 == 0 {
		return nil
	}
	return fmt.Errorf("%v is not even number", randomInt)
}

// ensure data consistency, db records with bind
func SyncDBRecordWithNsRR() (err error) {
	zones := []string{}
	err = database.DB.Model(&models.DnsZone{}).Pluck("zone", &zones).Error
	if err != nil {
		return
	}
	// name zone type content
	rrFmt := "%v/%v/%v/%v/%v"
	for _, zone := range zones {
		dnsRRsMap := make(map[string]struct{})
		dbRRsMap := make(map[string]struct{})

		dnsRRs, err := dnslib.IntranetRRsInZone(zone)
		if err != nil {
			return err
		}
		for _, record := range dnsRRs {
			str := fmt.Sprintf(rrFmt, record.RecordName, record.Zone, record.RecordType, record.RecordContent, record.RecordTtl)
			dnsRRsMap[str] = struct{}{}
		}

		dbRRs := []models.DnsRecord{}
		err = database.DB.Model(&models.DnsRecord{}).Where("zone = ?", zone).Find(&dbRRs).Error
		if err != nil {
			return err
		}
		for _, record := range dbRRs {
			str := fmt.Sprintf(rrFmt, record.RecordName, record.Zone, record.RecordType, record.RecordContent, record.RecordTtl)
			dbRRsMap[str] = struct{}{}
		}

		// diff
		for k := range dbRRsMap {
			if _, ok := dnsRRsMap[k]; ok {
				delete(dnsRRsMap, k)
				continue
			}
			array := strings.Split(k, "/")
			err = database.DB.Where("record_name = ?", array[0]).Where("zone = ?", array[1]).
				Where("record_type = ?", array[2]).Where("record_content = ?", array[3]).Delete(&models.DnsRecord{}).Error
			logrus.WithField("db rr", "del").Infof("%s %v", k, err)
		}
		toBeAddRecords := []models.DnsRecord{}
		for k := range dnsRRsMap {
			array := strings.Split(k, "/")
			ttl, _ := strconv.Atoi(array[4])
			toBeAddRecords = append(toBeAddRecords, models.DnsRecord{
				DnsRR: models.DnsRR{
					RecordName:    array[0],
					Zone:          array[1],
					RecordType:    array[2],
					RecordContent: array[3],
					RecordTtl:     ttl,
				},
				Creator: config.GlobalConfig.App.Name,
			})
		}
		err = database.DB.CreateInBatches(&toBeAddRecords, len(toBeAddRecords)).Error
		logrus.WithField("db rr", "add").Infof("%v %v", toBeAddRecords, err)
	}
	return
}

// dns probe, alarm is generated when the expected result is not met
func dnsProbe() (err error) {
	probes := []models.DnsProbe{}
	err = database.DB.Find(&probes).Error
	if err != nil {
		return
	}

	pool := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}
	for _, probe := range probes {
		wg.Add(1)
		pool <- struct{}{}
		go func(probe models.DnsProbe) {
			defer wg.Done()
			defer func() {
				<-pool
			}()
			answers := []string{}
			if probe.Intranet {
				rrs, err := dnslib.IntranetRRQueryAll(probe.RecordName, probe.Zone)
				if err != nil {
					probe.Result = models.ProbeQueryFailed
					database.DB.Save(&probe)
					return
				}
				for _, rr := range rrs {
					answers = append(answers, rr.RecordContent)
				}
			} else {
				rrs, err := dnslib.PublicDnsQueryRR(probe.RecordName)
				if err != nil {
					probe.Result = models.ProbeQueryFailed
					database.DB.Save(&probe)
					return
				}
				for _, rr := range rrs {
					answers = append(answers, rr.RecordContent)
				}
			}
			sort.Strings(answers)
			str := strings.Join(answers, ",")
			expectStr := strings.Join(probe.ExpectAnswer, ",")
			if str != expectStr {
				probe.Result = models.ProbeMisMatch
				database.DB.Save(&probe)
			}
		}(probe)
	}
	wg.Wait()
	return
}
