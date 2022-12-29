package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// SOA is ignored for now
type DNSDomain struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Note             string   `json:"note"`
	Status           string   `json:"status"`
	GeoIPEnabled     bool     `json:"geoip"`
	GTDEnabled       bool     `json:"gtd"`
	Nameservers      []string `json:"nameservers"`
	Tags             []string `json:"tags"`
	Template         int      `json:"template"`
	VanityNameserver []string `json:"vanityNameserver"`
	Contacts         []int    `json:"contacts"`
	CreatedAt        string   `json:"createdAt"`
	UpdatedAt        string   `json:"updatedAt"`
}

// GetDNSDomains returns active DNS domains in Constellix
func GetDNSDomains() ([]*DNSDomain, error) {
	logger.Println("Retrieving DNS domains...")
	endpoint, err := url.JoinPath(dnsRESTAPIBaseURL, "domains")
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
	var domains []*DNSDomain
	err = json.Unmarshal(resp.Data, &domains)
	if err != nil {
		return nil, err
	}
	return domains, nil
}
