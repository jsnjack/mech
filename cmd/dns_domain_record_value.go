package cmd

import (
	"fmt"
)

type DNSStandardValue struct {
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type DNSFailoverValue struct {
	Mode    string                  `json:"mode"`
	Enabled bool                    `json:"enabled"`
	Values  []*DNSFailoverItemValue `json:"values"`
}

type DNSFailoverItemValue struct {
	Enabled      bool   `json:"enabled"`
	Order        int    `json:"order"`
	SonarCheckID int    `json:"sonarCheckId"`
	Value        string `json:"value"`
}

type aliasDNSRecord DNSRecord

// populateDNSRecordValue populates the Value field of a DNSRecord based on the
// Mode field.
// TODO: be carefull with type casting, use similar to sonarCheckID everywhere
func populateDNSRecordValue(record interface{}) error {
	s, ok := record.(*DNSRecord)
	if !ok {
		return fmt.Errorf("unable to assert record to DNSRecord")
	}
	switch s.Mode {
	case "standard":
		m, ok := s.Value.([]interface{})
		if !ok {
			return fmt.Errorf("unable to parse value for standard mode, expected an array")
		}
		valueObj := make([]*DNSStandardValue, 0)
		for _, el := range m {
			elMap, ok := el.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for standard mode, expected an map")
			}
			valueEl := DNSStandardValue{
				Value:   elMap["value"].(string),
				Enabled: elMap["enabled"].(bool),
			}
			valueObj = append(valueObj, &valueEl)
		}
		s.Value = valueObj
	case "failover":
		valueObj := DNSFailoverValue{}
		m, ok := s.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to parse value for failover mode, expected an map")
		}
		valueObj.Mode = m["mode"].(string)
		valueObj.Enabled = m["enabled"].(bool)
		values := make([]*DNSFailoverItemValue, 0)
		for _, valueItem := range m["values"].([]interface{}) {
			valueItemMap, ok := valueItem.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for value of failover mode, expected an map")
			}
			valueItemObj := DNSFailoverItemValue{
				Enabled:      valueItemMap["enabled"].(bool),
				Order:        toInt(valueItemMap["order"]),
				Value:        valueItemMap["value"].(string),
				SonarCheckID: toInt(valueItemMap["sonarCheckId"]),
			}
			values = append(values, &valueItemObj)
		}
		valueObj.Values = values
		s.Value = &valueObj
	case "roundrobin-failover":
		// Round Robin Failover mode
	case "pools":
		// Pools mode
	default:
		return fmt.Errorf("unknown mode %q", s.Mode)
	}
	return nil
}

func toInt(i interface{}) int {
	switch v := i.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}
