package cmd

import (
	"testing"
)

func TestGetGeoproximityID_ValidStringInput(t *testing.T) {
	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	input := "@geoproximity:test"
	want := 1
	got, err := getGeoproximityID(input)
	if err != nil {
		t.Errorf("getGeoproximityID() error = %v, want error %v", err, false)
	}
	if got != want {
		t.Errorf("getGeoproximityID() = %v, want %v", got, want)
	}
}

func TestGetGeoproximityID_ValidStringInput2(t *testing.T) {
	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	input := "@geoproximity: test "
	want := 1
	got, err := getGeoproximityID(input)
	if err != nil {
		t.Errorf("getGeoproximityID() error = %v, want error %v", err, false)
	}
	if got != want {
		t.Errorf("getGeoproximityID() = %v, want %v", got, want)
	}
}

func TestGetGeoproximityID_ValidIntegerInput(t *testing.T) {
	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	input := 10
	want := 10
	got, err := getGeoproximityID(input)
	if err != nil {
		t.Errorf("getGeoproximityID() error = %v, want error %v", err, false)
	}
	if got != want {
		t.Errorf("getGeoproximityID() = %v, want %v", got, want)
	}
}

func TestGetGeoproximityID_InvalidStringInput(t *testing.T) {
	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	input := "test"
	expecxtedError := "invalid geoproximity value. Expected @geoproximity:<name> or int"
	_, err := getGeoproximityID(input)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err.Error() != expecxtedError {
		t.Errorf("getGeoproximityID() error = %v, want error %v", err, expecxtedError)
	}
}

func TestGetGeoproximityID_InvalidInputType(t *testing.T) {
	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	input := []int{1, 2, 3}
	expecxtedError := "invalid geoproximity value. Expected @geoproximity:<name> or int"
	_, err := getGeoproximityID(input)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err.Error() != expecxtedError {
		t.Errorf("getGeoproximityID() error = %v, want error %v", err, expecxtedError)
	}
}

func TestGetGeoproximityID_NotFound(t *testing.T) {
	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	input := "@geoproximity:unknown"
	expecxtedError := "unable to find geoproximity unknown"
	_, err := getGeoproximityID(input)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err.Error() != expecxtedError {
		t.Errorf("getGeoproximityID() error = %v, want error %v", err, expecxtedError)
	}
}

func TestPopulateDNSRecordGeoproximityForYAML_ValidInputString(t *testing.T) {
	record := &DNSRecord{
		GeoProximity: "@geoproximity:test",
	}

	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{
			{ID: 1, Name: "test"},
		}, nil
	}

	err := populateDNSRecordGeoproximityForYAML(record)
	if err != nil {
		t.Errorf("populateDNSRecordGeoproximityForYAML() error = %v, wantErr %v", err, false)
	}
	if record.GeoProximity != 1 {
		t.Errorf("populateDNSRecordGeoproximityForYAML() = %v, want %v", record.GeoProximity, 1)
	}
}

func TestPopulateDNSRecordGeoproximityForYAML_ValidInteger(t *testing.T) {
	record := &DNSRecord{
		GeoProximity: 1,
	}

	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{}, nil
	}

	err := populateDNSRecordGeoproximityForYAML(record)
	if err != nil {
		t.Errorf("populateDNSRecordGeoproximityForYAML() error = %v, wantErr %v", err, false)
	}
	if record.GeoProximity != 1 {
		t.Errorf("populateDNSRecordGeoproximityForYAML() = %v, want %v", record.GeoProximity, 1)
	}
}

func TestPopulateDNSRecordGeoproximityForYAML_Nil(t *testing.T) {
	record := &DNSRecord{
		GeoProximity: nil,
	}

	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{}, nil
	}

	err := populateDNSRecordGeoproximityForYAML(record)
	if err != nil {
		t.Errorf("populateDNSRecordGeoproximityForYAML() error = %v, wantErr %v", err, false)
	}
	if record.GeoProximity != nil {
		t.Errorf("populateDNSRecordGeoproximityForYAML() = %v, want %v", record.GeoProximity, nil)
	}
}

func TestPopulateDNSRecordGeoproximityForYAML_Invalid(t *testing.T) {
	record := &DNSRecord{
		GeoProximity: []int{1, 2, 3},
	}

	// Mock the GetGeoProximities function to control its behavior for testing
	oldGetGeoProximities := GetGeoProximities
	defer func() { GetGeoProximities = oldGetGeoProximities }()
	GetGeoProximities = func() ([]*GeoProximity, error) {
		return []*GeoProximity{}, nil
	}

	err := populateDNSRecordGeoproximityForYAML(record)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestPopulateDNSRecordGeoproximityForJSON_ValidObject(t *testing.T) {
	record := &DNSRecord{
		GeoProximity: map[string]interface{}{
			"id":   1,
			"name": "test",
		},
	}

	err := populateDNSRecordGeoproximityForJSON(record)
	if err != nil {
		t.Errorf("populateDNSRecordGeoproximityForJSON() error = %v, wantErr %v", err, false)
	}
	if record.GeoProximity != 1 {
		t.Errorf("populateDNSRecordGeoproximityForJSON() = %v, want %v", record.GeoProximity, 1)
	}
}

func TestPopulateDNSRecordGeoproximityForJSON_empty(t *testing.T) {
	record := &DNSRecord{}

	err := populateDNSRecordGeoproximityForJSON(record)
	if err != nil {
		t.Errorf("populateDNSRecordGeoproximityForJSON() error = %v, wantErr %v", err, false)
	}
	if record.GeoProximity != nil {
		t.Errorf("populateDNSRecordGeoproximityForJSON() = %v, want %v", record.GeoProximity, nil)
	}
}

func TestPopulateDNSRecordGeoproximityForJSON_invalid(t *testing.T) {
	record := &DNSRecord{
		GeoProximity: "hello",
	}

	err := populateDNSRecordGeoproximityForJSON(record)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
