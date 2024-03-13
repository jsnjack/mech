package cmd

import "fmt"

// In Constellix API, GET requests return an object, but POST and PATCH requests
// expect an integer. yaml configuration will also have an integer.

// populateDNSRecordIPFilterForJSON populates the IPFilter field from
// the JSON response from the API.
func populateDNSRecordIPFilterForJSON(record interface{}) error {
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

// populateDNSRecordIPFilterForYAML populates the IPFilter field from
// the local YAML configuration.
func populateDNSRecordIPFilterForYAML(record interface{}) error {
	s, ok := record.(*DNSRecord)
	if !ok {
		return fmt.Errorf("unable to assert record to DNSRecord")
	}
	if s.IPFilter == nil {
		return nil
	}
	switch s.IPFilter.(type) {
	case int, float64:
		s.IPFilter = toInt(s.IPFilter)
		return nil
	default:
		return fmt.Errorf("unable to parse value for ipfilter, expected an integer")
	}
	return nil
}
