package cmd

import "fmt"

// DNSIPFilter represents a ipfilter for a DNS record. In Constellix API, GET
// requests return an object, but POST and PATCH requests expect an integer.
// yaml configuration will also have an integer.
type DNSIPFilter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func populateDNSRecordIPFilter(record interface{}) error {
	s, ok := record.(*DNSRecord)
	if !ok {
		return fmt.Errorf("unable to assert record to DNSRecord")
	}
	if s.IPFilter == nil {
		return nil
	}
	elMap, ok := s.IPFilter.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to parse value for ipfilter, expected an object")
	}
	s.IPFilter = toInt(elMap["id"])
	return nil
}
