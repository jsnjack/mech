package cmd

import (
	"strings"
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

	dataExpected := "found change in immutable field: Host: 1.2.3.41.2.3.5"

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionError {
		t.Errorf("expected action '%v', got '%v'", ActionError, action)
		return
	}
	dataNoColor := stripColor(data)
	if dataExpected != dataNoColor {
		t.Errorf("expected data '%v', got '%v'", dataExpected, dataNoColor)
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
	action, data, err := Compare(&expected, &active)

	dataExpected1 := "Port: 44380"
	dateExpected2 := "Interval: NIGHTDAY"

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
	dataNoColor := stripColor(data)
	if strings.Contains(dataNoColor, dataExpected1) == false {
		t.Errorf("expected to contain '%v', got '%v'", dataExpected1, dataNoColor)
		return
	}
	if strings.Contains(dataNoColor, dateExpected2) == false {
		t.Errorf("expected to contain '%v', got '%v'", dateExpected2, dataNoColor)
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
	action, data, err := Compare(&expected, nil)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionCreate {
		t.Errorf("expected action '%v', got '%v'", ActionCreate, action)
		return
	}

	if data != "" {
		t.Errorf("expected data '%v', got '%v'", "", data)
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
