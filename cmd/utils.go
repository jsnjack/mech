package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Purple = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"

func colorAction(action ResourceAction) string {
	var start string
	switch action {
	case ActionCreate:
		start = Green
	case ActionDelete:
		start = Red
	case ActionUpate:
		start = Yellow
	case ActionOK:
		start = White
	}
	return start + string(action) + Reset
}

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

func makeAPIRequest(method string, url string, payload io.Reader, expectedStatusCode int) (respBody []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("x-cns-security-token", buildSecurityToken())
	req.Header.Add("Content-Type", "application/json")
	if rootVerbose {
		fmt.Printf("  requesting %s %s ...\n", method, url)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != expectedStatusCode {
		return body, fmt.Errorf("unexpected status code %d, want %d", res.StatusCode, expectedStatusCode)
	}
	return body, nil
}
