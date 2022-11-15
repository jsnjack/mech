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
