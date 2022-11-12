package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// SonarHTTPCheck represents sonar check record
// Example:
//
//	{
//	   "id": 81994,
//	   "name": "msi-prod-pdfeditor-backup",
//	   "host": "18.180.166.79",
//	   "port": 443,
//	   "protocolType": "HTTPS",
//	   "ipVersion": "IPV4",
//	   "fqdn": "msi-pdfeditor.surfly.jp",
//	   "path": "/blackbox/HealthCheck",
//	   "searchString": "",
//	   "connectionTimeout": 5,
//	   "expectedStatusCode": 200,
//	   "userAgent": "IE",
//	   "note": "",
//	   "runTraceroute": "DISABLED",
//	   "scheduleInterval": "NONE",
//	   "sslPolicy": "IGNORE",
//	   "userId": 300003895,
//	   "interval": "ONEMINUTE",
//	   "monitorIntervalPolicy": "PARALLEL",
//	   "checkSites": [
//	     15
//	   ],
//	   "notificationGroups": [
//	     41680,
//	     41675
//	   ],
//	   "scheduleId": 0,
//	   "notificationReportTimeout": 1440,
//	   "verificationPolicy": "SIMPLE"
//	 }
type SonarHTTPCheck struct {
	// Mandatory fields
	// TODO: Verify it!
	// Name is the unique identifier of Check
	Name          string `json:"name" yaml:"name"`
	Host          string `json:"host" yaml:"host"`
	IPVersion     string `json:"ipVersion" yaml:"ipVersion"`
	Port          int    `json:"port" yaml:"port"`
	ProtocolType  string `json:"protocolType" yaml:"protocolType"`
	Interval      string `json:"interval" yaml:"interval"`
	CheckSites    []int  `json:"checkSites" yaml:"checkSites"`
	RunTraceroute string `json:"runTraceroute" yaml:"runTraceroute"`

	FQDN                      string `json:"fqdn,omitempty" yaml:"fqdn"`
	Path                      string `json:"path,omitempty" yaml:"path"`
	SearchString              string `json:"searchString,omitempty" yaml:"searchString"`
	ConnectionTimeout         int    `json:"connectionTimeout,omitempty" yaml:"connectionTimeout"`
	ExpectedStatusCode        int    `json:"expectedStatusCode,omitempty" yaml:"expectedStatusCode"`
	UserAgent                 string `json:"userAgent,omitempty" yaml:"userAgent"`
	Note                      string `json:"note,omitempty" yaml:"note"`
	ScheduleInterval          string `json:"scheduleInterval,omitempty" yaml:"scheduleInterval"`
	SSLPolicy                 string `json:"sslPolicy,omitempty" yaml:"sslPolicy"`
	UserID                    int    `json:"userId,omitempty" yaml:"userId"`
	MonitorIntervalPolicy     string `json:"monitorIntervalPolicy,omitempty" yaml:"monitorIntervalPolicy"`
	NotificationGroups        []int  `json:"notificationGroups,omitempty" yaml:"notificationGroups"`
	ScheduleID                int    `json:"scheduleId,omitempty" yaml:"scheduleId"`
	NotificationReportTimeout int    `json:"notificationReportTimeout,omitempty" yaml:"notificationReportTimeout"`
	VerificationPolicy        string `json:"verificationPolicy,omitempty" yaml:"verificationPolicy"`
}

type ExpectedSonarHTTPCheck struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	SonarHTTPCheck
}

// UnmarshalJSON unmarshals the mesage and stores original fields
func (ex *ExpectedSonarHTTPCheck) UnmarshalJSON(b []byte) error {

	// Unmarshall data into SonarHTTPCheck struct
	var s SonarHTTPCheck
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	ex.SonarHTTPCheck = s

	// Save specified fields
	dm := make(map[string]interface{})
	err = json.Unmarshal(b, &dm)
	if err != nil {
		return err
	}

	definedFields := make([]string, len(dm))
	i := 0
	for k := range dm {
		definedFields[i] = k
		i++
	}
	ex.definedFieldsMap = getJSONTagsFromStruct(&ex.SonarHTTPCheck, definedFields...)

	return nil
}

// Compare compares ExpectedSonarHTTPCheck with active SonarHTTPCheck. Returns Status and data
func (e *ExpectedSonarHTTPCheck) Compare(activeResources *[]SonarHTTPCheck) (ResourceAction, []byte, error) {
	var active SonarHTTPCheck
	var found bool
	for _, el := range *activeResources {
		if el.Name == e.Name {
			active = el
			found = true
			break
		}
	}

	if !found {
		return ActionCreate, nil, nil
	}

	diffStructFields := make([]string, 0)

	expectedValue := reflect.ValueOf(e.SonarHTTPCheck)
	activeValue := reflect.ValueOf(active)
	for _, structFieldName := range e.definedFieldsMap {
		fmt.Printf("Inspecting field %s\n", structFieldName)
		fieldExpected := expectedValue.FieldByName(structFieldName)
		fieldActive := activeValue.FieldByName(structFieldName)
		// Compare field values
		if fieldExpected.Interface() != fieldActive.Interface() {
			diffStructFields = append(diffStructFields, structFieldName)
			fmt.Printf("Field %q got %q, want %q\n", structFieldName, fieldActive, fieldExpected)
		}
	}

	if len(diffStructFields) == 0 {
		return ActionOK, nil, nil
	}

	diffJSONFields := make([]string, 0)
	for k, v := range e.definedFieldsMap {
		if slices.Contains(diffStructFields, v) {
			diffJSONFields = append(diffJSONFields, k)
		}
	}

	dataBytes, err := toFilteredJSON(e.SonarHTTPCheck, diffJSONFields...)
	if err != nil {
		return "", nil, err
	}

	return ActionUpate, dataBytes, nil

}

// Verify verifies data for validity
func (e *ExpectedSonarHTTPCheck) Verify() error {
	return nil
}
