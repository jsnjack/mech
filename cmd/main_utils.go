package cmd

import (
	"encoding/json"
	"reflect"
	"strings"
)

func ToFilteredJSON(obj interface{}, includeFields ...string) ([]byte, error) {
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
		if contains(includeFields, key) {
			dataOut[key] = value
		}
	}

	dataOutBytes, err := json.Marshal(dataOut)
	if err != nil {
		return nil, err
	}
	return dataOutBytes, nil
}

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func findIndex[T comparable](collection []T, el T) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}

// StructTagsToFieldNames returns struct field names from their tags
func StructTagsToFieldNames(obj interface{}, tags ...string) []string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	fields := reflect.VisibleFields(t)
	var filedsNames []string
OUTER:
	for _, tag := range tags {
		for _, f := range fields {
			val, ok := f.Tag.Lookup("json")
			if ok {
				tagName := strings.Split(val, ",")[0]
				if tagName == tag {
					filedsNames = append(filedsNames, f.Name)
					continue OUTER
				}
			}
		}
	}
	return filedsNames
}
