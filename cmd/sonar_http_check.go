package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	yaml "gopkg.in/yaml.v3"
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
	// Name should be the unique identifier of Check
	ID                        int    `json:"id"`
	Name                      string `json:"name" yaml:"name"`
	Host                      string `json:"host" yaml:"host"`
	IPVersion                 string `json:"ipVersion" yaml:"ipVersion"`
	Port                      int    `json:"port" yaml:"port"`
	ProtocolType              string `json:"protocolType" yaml:"protocolType"`
	Interval                  string `json:"interval" yaml:"interval"`
	CheckSites                []int  `json:"checkSites" yaml:"checkSites"`
	RunTraceroute             string `json:"runTraceroute" yaml:"runTraceroute"`
	FQDN                      string `json:"fqdn" yaml:"fqdn"`
	Path                      string `json:"path" yaml:"path"`
	SearchString              string `json:"searchString" yaml:"searchString"`
	ConnectionTimeout         int    `json:"connectionTimeout" yaml:"connectionTimeout"`
	ExpectedStatusCode        int    `json:"expectedStatusCode" yaml:"expectedStatusCode"`
	UserAgent                 string `json:"userAgent" yaml:"userAgent"`
	Note                      string `json:"note" yaml:"note"`
	ScheduleInterval          string `json:"scheduleInterval" yaml:"scheduleInterval"`
	SSLPolicy                 string `json:"sslPolicy" yaml:"sslPolicy"`
	UserID                    int    `json:"userId" yaml:"userId"`
	MonitorIntervalPolicy     string `json:"monitorIntervalPolicy" yaml:"monitorIntervalPolicy"`
	NotificationGroups        []int  `json:"notificationGroups" yaml:"notificationGroups"`
	ScheduleID                int    `json:"scheduleId" yaml:"scheduleId"`
	NotificationReportTimeout int    `json:"notificationReportTimeout" yaml:"notificationReportTimeout"`
	VerificationPolicy        string `json:"verificationPolicy" yaml:"verificationPolicy"`
}

type ExpectedSonarHTTPCheck struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	SonarHTTPCheck
	// This fields must be present when running update action
	mandatoryStructFields []string
}

// UnmarshalYAML unmarshals the mesage and stores original fields
func (ex *ExpectedSonarHTTPCheck) UnmarshalYAML(value *yaml.Node) error {

	// Unmarshall data into SonarHTTPCheck struct
	var s SonarHTTPCheck
	err := value.Decode(&s)
	if err != nil {
		return err
	}
	ex.SonarHTTPCheck = s

	// Save specified fields
	dm := make(map[string]interface{})
	err = value.Decode(&dm)
	if err != nil {
		return err
	}

	definedFields := make([]string, len(dm))
	i := 0
	for k := range dm {
		definedFields[i] = k
		i++
	}
	ex.definedFieldsMap = getFieldNamesMap(&ex.SonarHTTPCheck, "yaml", definedFields...)
	ex.mandatoryStructFields = []string{"Name", "ProtocolType", "Port"}

	return nil
}

// Compare compares ExpectedSonarHTTPCheck with active SonarHTTPCheck. Returns Status and data
func (e *ExpectedSonarHTTPCheck) Compare(activeResources *[]SonarHTTPCheck) (ResourceAction, []byte, error) {
	var active SonarHTTPCheck
	el, found := e.GetActive(activeResources)
	// el is a pointer, but reflect code expects real object
	if found {
		active = *el
	}

	var action ResourceAction

	diffStructFields := make([]string, 0)
	if !found {
		action = ActionCreate
		diffStructFields = maps.Values(e.definedFieldsMap)
	} else {
		action = ActionOK
		expectedValue := reflect.ValueOf(e.SonarHTTPCheck)
		activeValue := reflect.ValueOf(active)
		for _, structFieldName := range e.definedFieldsMap {
			fieldExpected := expectedValue.FieldByName(structFieldName)
			fieldActive := activeValue.FieldByName(structFieldName)
			// Compare field values
			if !reflect.DeepEqual(fieldExpected.Interface(), fieldActive.Interface()) {
				action = ActionUpate
				// Sonar blocks UPDATE requests with these parameters
				diffStructFields = removeItems(maps.Values(e.definedFieldsMap), "Host", "IPVersion")
				break
			}
		}

		if action == ActionOK {
			return action, nil, nil
		}
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

	return action, dataBytes, nil
}

func (e *ExpectedSonarHTTPCheck) GetActive(activeResources *[]SonarHTTPCheck) (*SonarHTTPCheck, bool) {
	for _, el := range *activeResources {
		if el.Name == e.Name {
			return &el, true
		}
	}
	return nil, false
}

// GetSonarHTTPChecks returns active Sonar Checks
func GetSonarHTTPChecks() (*[]SonarHTTPCheck, error) {
	// Fetch HTTP checks
	fmt.Println("Retrieving Sonar HTTP Checks...")
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http")
	if err != nil {
		return nil, err
	}
	data, err := makeAPIRequest("GET", endpoint, nil, 200)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sonar HTTP checks: %s", err)
	}

	checks := make([]SonarHTTPCheck, 0)
	err = json.Unmarshal(data, &checks)
	if err != nil {
		return nil, err
	}
	return &checks, nil
}

func CreateSonarHTTPCheck(payload []byte) error {
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http")
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeAPIRequest("POST", endpoint, payloadReader, 201)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to create Sonar HTTP checks: %s", err)
	}
	return nil
}

func UpdateSonarHTTPCheck(payload []byte, id int) error {
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http", fmt.Sprint(id))
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeAPIRequest("PUT", endpoint, payloadReader, 200)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to update Sonar HTTP checks: %s", err)
	}
	return nil
}

func DeleteSonarHTTPCheck(payload []byte, id int) error {
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http", fmt.Sprint(id))
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeAPIRequest("DELETE", endpoint, payloadReader, 202)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to delete Sonar HTTP checks: %s", err)
	}
	return nil
}
