package pkg

import (
	"net"
	"strings"
)

type EmailData struct {
	HasMX    bool
	HasSPF   bool
	HasDMARC bool

	SpfRecord   string
	DmarcRecord string
}

func CheckDomain(domain string) (EmailData, error) {
	data := EmailData{}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return data, err
	}
	if len(mxRecords) > 0 {
		data.HasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return data, err
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			data.HasSPF = true
			data.SpfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return data, err
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			data.HasDMARC = true
			data.DmarcRecord = record
			break
		}
	}

	return data, nil
}
