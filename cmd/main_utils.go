package cmd

import (
	"encoding/json"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

// toFilteredJSON mashals struct into JSON bytes which contain only specified fields
func toFilteredJSON(obj interface{}, includeFields ...string) ([]byte, error) {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Convert obj to map to simplify iteration
	dataIn := map[string]interface{}{}
	json.Unmarshal(objBytes, &dataIn)

	dataOut := map[string]interface{}{}

	// Create a new data obj which contains only fields which need to be included
	for key, value := range dataIn {
		if slices.Contains(includeFields, key) {
			dataOut[key] = value
		}
	}

	dataOutBytes, err := json.Marshal(dataOut)
	if err != nil {
		return nil, err
	}
	return dataOutBytes, nil
}

// getFieldNamesMap returns struct field names from their tags
func getFieldNamesMap(obj interface{}, tagType string, tags ...string) map[string]string {
	res := make(map[string]string)
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	fields := reflect.VisibleFields(t)
OUTER:
	for _, tag := range tags {
		for _, f := range fields {
			val, ok := f.Tag.Lookup(tagType)
			if ok {
				tagName := strings.Split(val, ",")[0]
				if tagName == tag {
					res[tag] = f.Name
					continue OUTER
				}
			}
		}
	}
	return res
}
