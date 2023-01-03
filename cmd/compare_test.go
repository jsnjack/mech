package cmd

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func Test_Compare_http_checks_diff_port(t *testing.T) {
	expectedStr := `
name: prod
port: 80
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", Port: 443}
	action, _, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
}

func Test_Compare_http_checks_diff_eq(t *testing.T) {
	expectedStr := `
name: prod
port: 80
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", Port: 80}
	action, _, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionOK {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}
}

func Test_Compare_http_checks_diff_slice_neq(t *testing.T) {
	expectedStr := `
name: prod
checkSites: [1,2]
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", CheckSites: []int{2}}
	action, _, err := Compare(&expected, &active)

	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}

	if err != nil {
		t.Errorf("expected no err, got '%v'", err)
		return
	}
}

func Test_Compare_http_checks_diff_slice_neq_order(t *testing.T) {
	expectedStr := `
name: prod
checkSites: [1,2]
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", CheckSites: []int{2, 1}}
	action, _, err := Compare(&expected, &active)

	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}

	if err != nil {
		t.Errorf("expected no err, got '%v'", err)
		return
	}
}

func Test_Compare_http_checks_diff_port_extra_active_field(t *testing.T) {
	expectedStr := `
name: prod
port: 80
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", Port: 443, Interval: "HALFDAY"}
	action, _, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
}

func Test_Compare_http_checks_diff_port_unknown_expected_field(t *testing.T) {
	expectedStr := `
name: prod
port: 80
unknown: 20
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", Port: 443, Interval: "HALFDAY"}
	action, _, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
}

func Test_Compare_http_checks_immunable_field(t *testing.T) {
	expectedStr := `
name: prod
port: 80
unknown: 20
host: 1.2.3.5
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", Port: 443, Host: "1.2.3.4"}
	action, data, err := Compare(&expected, &active)

	errorExpected := "found change in immutable field: Host"
	if action != ActionError {
		t.Errorf("expected action '%v', got '%v'", ActionError, action)
		return
	}

	if err == nil {
		t.Errorf("expected err, got nil")
		return
	}

	if err.Error() != errorExpected {
		t.Errorf("expected err '%v', got '%v'", errorExpected, err.Error())
		return
	}

	if len(data) != 0 {
		t.Errorf("expected no data, got '%v'", data)
		return
	}
}

func Test_Compare_http_checks_multiple_fields(t *testing.T) {
	expectedStr := `
name: prod
port: 80
interval: DAY
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	active := SonarHTTPCheck{Name: "prod", Port: 443, Interval: "NIGHT"}
	action, diffs, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
	portDiff := getFieldDiff(diffs, "Port")
	if portDiff == nil {
		t.Errorf("expected port diff, got nil")
		return
	}
	if portDiff.OldValue != "443" && portDiff.NewValue != "80" {
		t.Errorf("unexpected diff %s, want old value 80, vew value 443", portDiff.String())
		return
	}
}

func Test_Compare_http_checks_create(t *testing.T) {
	expectedStr := `
name: prod
port: 80
`
	var expected ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	action, diffs, err := Compare(&expected, nil)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionCreate {
		t.Errorf("expected action '%v', got '%v'", ActionCreate, action)
		return
	}

	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %v", diffs)
		return
	}
}

func Test_ToResourceMatcher(t *testing.T) {
	var collection []*ExpectedSonarHTTPCheck
	collection = append(collection, &ExpectedSonarHTTPCheck{})
	result := toResourceMatcher(collection)
	if len(result) != 1 {
		t.Errorf("expected 1, got %v", len(result))
		return
	}
}

func Test_Compare_dns_record_standard_value(t *testing.T) {
	expectedStr := `
name: abc
type: A
ttl: 60
mode: standard
region: default
enabled: true
value:
  - value: 1.1.1.1
    enabled: true
`
	var expected ExpectedDNSRecord
	err := yaml.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	expectedValue := make([]*DNSStandardValue, 0)
	expectedValue = append(expectedValue, &DNSStandardValue{Value: "8.8.8.8", Enabled: true})
	active := DNSRecord{
		Name:    "abc",
		Type:    "A",
		TTL:     60,
		Mode:    "standard",
		Region:  "default",
		Enabled: true,
		Value:   expectedValue,
	}
	action, _, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
}
