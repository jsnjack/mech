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
