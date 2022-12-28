package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type DNSDomain struct {
	ID int `json:"id"`
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
	fmt.Println(string(data))
	resp := DNSv4Response{}
	err = json.Unmarshal(data, &resp)
	fmt.Printf("%+v\n", resp)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
