package cmd

import (
	"testing"
)

func TestStructTagsToFieldNames_empty_struct(t *testing.T) {
	a := struct {
		name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	result := StructTagsToFieldNames(a, tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestStructTagsToFieldNames_empty_pointer(t *testing.T) {
	a := struct {
		name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	result := StructTagsToFieldNames(&a, tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestStructTagsToFieldNames_struct_public_empty(t *testing.T) {
	a := struct {
		Name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	result := StructTagsToFieldNames(a, tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestStructTagsToFieldNames_struct_public_no_match(t *testing.T) {
	a := struct {
		Name string
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := StructTagsToFieldNames(a, tags...)
	if len(result) != 0 {
		t.Errorf("expected empty, got %s\n", result)
	}
}

func TestStructTagsToFieldNames_struct_public_match1(t *testing.T) {
	a := struct {
		Name string `json:"name"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := StructTagsToFieldNames(a, tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	if result[0] != "Name" {
		t.Errorf("expected Name , got %s\n", result[0])
		return
	}
}

func TestStructTagsToFieldNames_struct_public_match2(t *testing.T) {
	a := struct {
		Name string `json:"name,omitempty"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := StructTagsToFieldNames(a, tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	if result[0] != "Name" {
		t.Errorf("expected Name , got %s\n", result[0])
		return
	}
}

func TestStructTagsToFieldNames_pointer_public_match(t *testing.T) {
	a := struct {
		Name string `json:"name,omitempty"`
		age  int
	}{
		"Joe", 10,
	}
	var tags []string
	tags = append(tags, "name")
	result := StructTagsToFieldNames(&a, tags...)
	if len(result) != 1 {
		t.Errorf("expected one element, got %d\n", len(result))
		return
	}
	if result[0] != "Name" {
		t.Errorf("expected Name , got %s\n", result[0])
		return
	}
}
