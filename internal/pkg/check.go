package pkg

import (
	"errors"
	"net"
	"strings"
)

type EmailData struct {
	hasMX    bool
	hasSPF   bool
	hasDMARC bool

	spfRecord   string
	dmarcRecord string
}

func CheckDomain(domain string) (EmailData, error) {
	data := EmailData{}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return data, err
	}
	if len(mxRecords) > 0 {
		data.hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return data, errors.New("Error: " + err.Error())
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			data.hasSPF = true
			data.spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return data, err
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			data.hasDMARC = true
			data.dmarcRecord = record
			break
		}
	}

	return data, nil
}
