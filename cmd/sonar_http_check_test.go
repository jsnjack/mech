package cmd

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestExpectedSonarHTTPCheck_UnmarshalYAML(t *testing.T) {
	data := `
name: prod
port: 80
alien: yes
`
	var obj ExpectedSonarHTTPCheck
	err := yaml.Unmarshal([]byte(data), &obj)
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

func TestExpectedSonarHTTPCheck_Validate_no_mandatory(t *testing.T) {
	data := `
name: prod
port: 80
`
	var obj ExpectedSonarHTTPCheck
	// Stub mandatory fields
	obj.mandatoryFields = []string{"name", "port", "host"}
	err := yaml.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	err = obj.Validate()
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	expected := "prod: mandatory field \"host\" is not defined"
	if err != nil && err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
		return
	}
}
