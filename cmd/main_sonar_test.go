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

func TestExpectedSonarHTTPCheck_Unmarshal(t *testing.T) {
	data := `{"name": "prod", "port": 80}`
	var obj ExpectedSonarHTTPCheck
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if len(obj.specifiedFields) != 2 {
		t.Errorf("wrong specified fields: %s, want name and port", obj.specifiedFields)
	}
}
