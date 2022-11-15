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
	action, data, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	expectedData := `{"name":"prod","port":80}`
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
	if string(data) != expectedData {
		t.Errorf("expected %v, got %v", expectedData, string(data))
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
	action, data, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	if action != ActionOK {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}
	if data != nil {
		t.Errorf("expected no data, got '%v'", data)
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
	action, data, err := Compare(&expected, &active)

	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}

	expectedData := `{"checkSites":[1,2],"name":"prod"}`
	if string(data) != expectedData {
		t.Errorf("expected no data, got '%v'", string(data))
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
	action, data, err := Compare(&expected, &active)

	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}

	expectedData := `{"checkSites":[1,2],"name":"prod"}`
	if string(data) != expectedData {
		t.Errorf("expected no data, got '%v'", string(data))
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
	action, data, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	expectedData := `{"name":"prod","port":80}`
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
	if string(data) != expectedData {
		t.Errorf("expected %v, got %v", expectedData, string(data))
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
	action, data, err := Compare(&expected, &active)

	if err != nil {
		t.Error(err)
		return
	}
	expectedData := `{"name":"prod","port":80}`
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
	if string(data) != expectedData {
		t.Errorf("expected %v, got %v", expectedData, string(data))
		return
	}
}
