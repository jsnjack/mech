package cmd

import (
	"testing"
)

func TestGetJSONTagsFromStruct_empty_struct(t *testing.T) {
	a := struct {
		name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	result := getFieldNamesMap(a, "yaml", tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestGetJSONTagsFromStruct_empty_pointer(t *testing.T) {
	a := struct {
		name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	result := getFieldNamesMap(&a, "yaml", tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestGetJSONTagsFromStruct_struct_public_empty(t *testing.T) {
	a := struct {
		Name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	result := getFieldNamesMap(a, "yaml", tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestGetJSONTagsFromStruct_struct_public_no_match(t *testing.T) {
	a := struct {
		Name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := getFieldNamesMap(a, "yaml", tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestGetJSONTagsFromStruct_struct_public_match1(t *testing.T) {
	a := struct {
		Name string `yaml:"name"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := getFieldNamesMap(a, "yaml", tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	expectedKey := "name"
	expectedValue := "Name"
	v, ok := result[expectedKey]
	if !ok {
		t.Errorf("expected to have key %q, got %q", expectedKey, result)
		return
	}
	if v != expectedValue {
		t.Errorf("expected %q to be mapped to %q, got %q", expectedKey, expectedValue, v)
		return
	}
}

func TestGetJSONTagsFromStruct_struct_public_match2(t *testing.T) {
	a := struct {
		Name string `yaml:"name,omitempty"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := getFieldNamesMap(a, "yaml", tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	expectedKey := "name"
	expectedValue := "Name"
	v, ok := result[expectedKey]
	if !ok {
		t.Errorf("expected to have key %q, got %q", expectedKey, result)
		return
	}
	if v != expectedValue {
		t.Errorf("expected %q to be mapped to %q, got %q", expectedKey, expectedValue, v)
		return
	}
}

func TestGetJSONTagsFromStruct_pointer_public_match(t *testing.T) {
	a := struct {
		Name string `json:"name,omitempty"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := getFieldNamesMap(&a, "json", tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	expectedKey := "name"
	expectedValue := "Name"
	v, ok := result[expectedKey]
	if !ok {
		t.Errorf("expected to have key %q, got %q", expectedKey, result)
		return
	}
	if v != expectedValue {
		t.Errorf("expected %q to be mapped to %q, got %q", expectedKey, expectedValue, v)
		return
	}
}

func TestGetJSONTagsFromStruct_with_yaml(t *testing.T) {
	a := struct {
		Name string `json:"name,omitempty" yaml:"yaname"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := getFieldNamesMap(&a, "json", tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	expectedKey := "name"
	expectedValue := "Name"
	v, ok := result[expectedKey]
	if !ok {
		t.Errorf("expected to have key %q, got %q", expectedKey, result)
		return
	}
	if v != expectedValue {
		t.Errorf("expected %q to be mapped to %q, got %q", expectedKey, expectedValue, v)
		return
	}
}
