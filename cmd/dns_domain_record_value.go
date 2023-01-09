package cmd

import (
	"fmt"
	"strings"
)

type DNSStandardItemValue struct {
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

type DNSMXStandardItemValue struct {
	Server   string `json:"server"`
	Priority int    `json:"priority"`
	Enabled  bool   `json:"enabled"`
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
	switch s.Type {
	case "A":
		switch s.Mode {
		case "standard":
			m, ok := s.Value.([]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for standard mode, expected an array")
			}
			valueObj := make([]*DNSStandardItemValue, 0)
			for _, el := range m {
				elMap, ok := el.(map[string]interface{})
				if !ok {
					return fmt.Errorf("unable to parse value for standard mode, expected an map")
				}
				valueEl := DNSStandardItemValue{
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
				sonarCheckID, err := getSonarCheckID(valueItemMap["sonarCheckId"])
				if err != nil {
					return err
				}
				valueItemObj := DNSFailoverItemValue{
					Enabled:      valueItemMap["enabled"].(bool),
					Order:        toInt(valueItemMap["order"]),
					Value:        valueItemMap["value"].(string),
					SonarCheckID: sonarCheckID,
				}
				values = append(values, &valueItemObj)
			}
			valueObj.Values = values
			s.Value = &valueObj
		case "roundrobin-failover":
			m, ok := s.Value.([]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for roundrobin-failover mode, expected an array")
			}
			valueObj := make([]*DNSFailoverItemValue, 0)
			for _, el := range m {
				elMap, ok := el.(map[string]interface{})
				if !ok {
					return fmt.Errorf("unable to parse value for roundrobin-failover mode, expected an map")
				}
				sonarCheckID, err := getSonarCheckID(elMap["sonarCheckId"])
				if err != nil {
					return err
				}
				valueEl := DNSFailoverItemValue{
					Value:        elMap["value"].(string),
					Enabled:      elMap["enabled"].(bool),
					Order:        toInt(elMap["order"]),
					SonarCheckID: sonarCheckID,
				}
				valueObj = append(valueObj, &valueEl)
			}
			s.Value = valueObj
		case "pools":
			m, ok := s.Value.([]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for pools mode, expected an array")
			}
			valueObj := make([]int, 0)
			for _, el := range m {
				valueObj = append(valueObj, toInt(el))
			}
			s.Value = valueObj
		default:
			return fmt.Errorf("unknown mode %q", s.Mode)
		}
	case "MX":
		if s.Mode != "standard" {
			return fmt.Errorf("unsupported mode %q for MX record", s.Mode)
		}
		m, ok := s.Value.([]interface{})
		if !ok {
			return fmt.Errorf("unable to parse value for MX record in standard mode, expected an array")
		}
		valueObj := make([]*DNSMXStandardItemValue, 0)
		for _, el := range m {
			elMap, ok := el.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for standard mode, expected an map")
			}
			valueEl := DNSMXStandardItemValue{
				Server:   elMap["server"].(string),
				Priority: toInt(elMap["priority"]),
				Enabled:  elMap["enabled"].(bool),
			}
			valueObj = append(valueObj, &valueEl)
		}
		s.Value = valueObj
	case "TXT":
		if s.Mode != "standard" {
			return fmt.Errorf("unsupported mode %q for TXT record", s.Mode)
		}
		m, ok := s.Value.([]interface{})
		if !ok {
			return fmt.Errorf("unable to parse value for TXT record in standard mode, expected an array")
		}
		valueObj := make([]*DNSStandardItemValue, 0)
		for _, el := range m {
			elMap, ok := el.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unable to parse value for TXT record in standard mode, expected an map")
			}
			valueEl := DNSStandardItemValue{
				Value:   elMap["value"].(string),
				Enabled: elMap["enabled"].(bool),
			}
			valueObj = append(valueObj, &valueEl)
			s.Value = valueObj
		}
	default:
		return fmt.Errorf("unsupported record type %q", s.Type)
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

func getSonarCheckID(i interface{}) (int, error) {
	switch v := i.(type) {
	case string:
		checkType, checkName, err := parseSonarCheckID(v)
		if err != nil {
			return 0, err
		}
		switch checkType {
		case "http":
			checks, err := GetSonarHTTPChecks()
			if err != nil {
				return 0, err
			}
			for _, check := range checks {
				if check.Name == checkName {
					return check.ID, nil
				}
			}
		case "tcp":
			checks, err := GetSonarTCPChecks()
			if err != nil {
				return 0, err
			}
			for _, check := range checks {
				if check.Name == checkName {
					return check.ID, nil
				}
			}
		}
	}
	return toInt(i), nil
}

// parseSonarCheckID parses a sonar check ID from a string. It assumes that the string
// will start with a @, followed by code word 'sonar' with specified check type and the
// name of the check itself
func parseSonarCheckID(s string) (string, string, error) {
	if !strings.HasPrefix(s, "@sonar,") {
		return "", "", fmt.Errorf("invalid sonar check ID. Expected @sonar,<check_type>:<check_name> or int")
	}
	s = strings.TrimPrefix(s, "@sonar,")
	split := strings.Split(s, ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("invalid sonar check ID. Expected @sonar,<check_type>:<check_name> or int")
	}
	return split[0], split[1], nil
}
