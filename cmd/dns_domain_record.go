package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Missing fields: lastValues, skipLookup, contacts
type DNSRecord struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	TTL            int              `json:"ttl"`
	Mode           string           `json:"mode"`
	Region         string           `json:"region"`
	IPFilter       *DNSIPFilter     `json:"ipfilter"`
	IPFilterIPDrop bool             `json:"ipfilteripDrop"`
	GeoFailover    bool             `json:"geoFailover"`
	GeoProximity   *DNSGeoProximity `json:"geoProximity"`
	Enabled        bool             `json:"enabled"`
	Value          []*DNSValue      `json:"value"`
	Notes          string           `json:"notes"`
}

type DNSIPFilter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DNSGeoProximity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DNSValue struct {
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

// GetDNSRecords retrieves domain's DNS records
func GetDNSRecords(id int) ([]*DNSRecord, error) {
	logger.Printf("Retrieving DNS records for domain %d...\n", id)
	endpoint, err := url.JoinPath(dnsRESTAPIBaseURL, "domains", fmt.Sprintf("%d", id), "records")
	if err != nil {
		return nil, err
	}
	data, err := makeSimpleAPIRequest("GET", endpoint, nil, 200)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve DNS domains list: %s", err)
	}
	resp := DNSv4Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	var records []*DNSRecord
	err = json.Unmarshal(resp.Data, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}
