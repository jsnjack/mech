package cmd

import (
	"encoding/json"
	"testing"
)

func TestExpectedSonarHTTPCheck_Compare_not_found(t *testing.T) {
	expected := ExpectedSonarHTTPCheck{}
	expected.SonarHTTPCheck = SonarHTTPCheck{Name: "prod"}
	activeList := make([]SonarHTTPCheck, 0)
	action, data, _ := expected.Compare(&activeList)
	if action != ActionCreate {
		t.Errorf("expected action '%v', got '%v'", ActionCreate, action)
		return
	}
	if data != nil {
		t.Errorf("expected nil data, got %s", string(data))
		return
	}
}

func TestExpectedSonarHTTPCheck_Compare_different(t *testing.T) {
	expectedStr := `{"name": "prod", "port": 80}`
	var expected ExpectedSonarHTTPCheck
	err := json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	activeList := make([]SonarHTTPCheck, 0)
	activeList = append(activeList, SonarHTTPCheck{Name: "prod", Port: 443})
	action, data, _ := expected.Compare(&activeList)
	expectedData := `{"port":80}`
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionUpate, action)
		return
	}
	if string(data) != expectedData {
		t.Errorf("expected %v, got %v", expectedData, string(data))
		return
	}
}

func TestExpectedSonarHTTPCheck_Compare_different_slice_eq(t *testing.T) {
	expectedStr := `{"name": "prod", "checkSites": [1,2]}`
	var expected ExpectedSonarHTTPCheck
	err := json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	activeList := make([]SonarHTTPCheck, 0)
	activeList = append(activeList, SonarHTTPCheck{Name: "prod", CheckSites: []int{1, 2}})
	action, data, err := expected.Compare(&activeList)
	if action != ActionOK {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}
	if data != nil {
		t.Errorf("expected no data, got '%v'", data)
		return
	}
	if err != nil {
		t.Errorf("expected no err, got '%v'", err)
		return
	}

}

func TestExpectedSonarHTTPCheck_Compare_different_slice_neq(t *testing.T) {
	expectedStr := `{"name": "prod", "checkSites": [1,2]}`
	var expected ExpectedSonarHTTPCheck
	err := json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	activeList := make([]SonarHTTPCheck, 0)
	activeList = append(activeList, SonarHTTPCheck{Name: "prod", CheckSites: []int{2}})
	action, data, err := expected.Compare(&activeList)
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}

	expectedData := `{"checkSites":[1,2]}`
	if string(data) != expectedData {
		t.Errorf("expected no data, got '%v'", string(data))
		return
	}
	if err != nil {
		t.Errorf("expected no err, got '%v'", err)
		return
	}

}

func TestExpectedSonarHTTPCheck_Compare_different_slice_neq_order(t *testing.T) {
	expectedStr := `{"name": "prod", "checkSites": [1,2]}`
	var expected ExpectedSonarHTTPCheck
	err := json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		t.Error(err)
		return
	}
	activeList := make([]SonarHTTPCheck, 0)
	activeList = append(activeList, SonarHTTPCheck{Name: "prod", CheckSites: []int{2, 1}, IPVersion: "ipv4"})
	action, data, err := expected.Compare(&activeList)
	if action != ActionUpate {
		t.Errorf("expected action '%v', got '%v'", ActionOK, action)
		return
	}

	expectedData := `{"checkSites":[1,2]}`
	if string(data) != expectedData {
		t.Errorf("expected no data, got '%v'", string(data))
		return
	}
	if err != nil {
		t.Errorf("expected no err, got '%v'", err)
		return
	}

}

func TestExpectedSonarHTTPCheck_UnmarshalJSON(t *testing.T) {
	data := `{"name": "prod", "port": 80, "alien": true}`
	var obj ExpectedSonarHTTPCheck
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if len(obj.definedFieldsMap) != 2 {
		t.Errorf("wrong length: got %d, want %d", len(obj.definedFieldsMap), 2)
		return
	}
	if obj.definedFieldsMap["name"] != "Name" {
		t.Errorf("expected %q to be mapped to %q, got %q", "name", "Name", obj.definedFieldsMap["name"])
		return
	}
	if obj.definedFieldsMap["port"] != "Port" {
		t.Errorf("expected %q to be mapped to %q, got %q", "port", "Port", obj.definedFieldsMap["port"])
		return
	}
}
