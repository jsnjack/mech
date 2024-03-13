package cmd

import (
	"fmt"
	"strings"
)

// In Constellix API, GET requests returns an object, but POST and PATCH requests
// expect an integer. yaml configuration will have expandable value.

// populateDNSRecordGeoproximityForJSON populates the GeoProximity field from
// the JSON response from the API.
func populateDNSRecordGeoproximityForJSON(record interface{}) error {
	s, ok := record.(*DNSRecord)
	if !ok {
		return fmt.Errorf("unable to assert record to DNSRecord")
	}
	if s.GeoProximity == nil {
		return nil
	}
	elMap, ok := s.GeoProximity.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to parse value for geoproximity, expected an object")
	}
	s.GeoProximity = toInt(elMap["id"])
	return nil
}

// populateDNSRecordGeoproximityForYAML populates the GeoProximity field from
// the local YAML configuration.
func populateDNSRecordGeoproximityForYAML(record interface{}) error {
	s, ok := record.(*DNSRecord)
	if !ok {
		return fmt.Errorf("unable to assert record to DNSRecord")
	}
	if s.GeoProximity == nil {
		return nil
	}
	gpID, err := getGeoproximityID(s.GeoProximity)
	if err != nil {
		return err
	}
	s.GeoProximity = gpID
	return nil
}

// getGeoproximityID returns the ID of the geoproximity object. It supports both
// an integer and a string `@georpximity:Name`.
func getGeoproximityID(gp interface{}) (int, error) {
	switch v := gp.(type) {
	case string:
		if !strings.HasPrefix(v, "@geoproximity:") {
			return 0, fmt.Errorf("invalid geoproximity value. Expected @geoproximity:<name> or int")
		}
		name := strings.TrimPrefix(v, "@geoproximity:")
		name = strings.TrimSpace(name)
		proximities, err := GetGeoProximities()
		if err != nil {
			return 0, err
		}
		for _, p := range proximities {
			if p.Name == name {
				return p.ID, nil
			}
		}
		return 0, fmt.Errorf("unable to find geoproximity %s", name)
	case int, float64:
		return toInt(gp), nil
	default:
		return 0, fmt.Errorf("invalid geoproximity value. Expected @geoproximity:<name> or int")
	}
}
