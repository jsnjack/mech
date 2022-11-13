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

// Example:
//
//	{
//	  "id": 83645,
//	  "name": "tcp-test",
//	  "host": "159.69.18.28",
//	  "port": 443,
//	  "ipVersion": "IPV4",
//	  "stringToSend": "",
//	  "stringToReceive": "",
//	  "note": "",
//	  "runTraceroute": "DISABLED",
//	  "userId": 300003895,
//	  "interval": "THIRTYSECONDS",
//	  "monitorIntervalPolicy": "PARALLEL",
//	  "checkSites": [
//	    4
//	  ],
//	  "notificationGroups": [],
//	  "scheduleId": 0,
//	  "notificationReportTimeout": 1440,
//	  "verificationPolicy": "SIMPLE"
//	}
type SonarTCPCheck struct {
	ID                        int    `json:"id"`
	Name                      string `json:"name" yaml:"name"`
	Host                      string `json:"host" yaml:"host"`
	IPVersion                 string `json:"ipVersion" yaml:"ipVersion"`
	Port                      int    `json:"port" yaml:"port"`
	Interval                  string `json:"interval" yaml:"interval"`
	CheckSites                []int  `json:"checkSites" yaml:"checkSites"`
	RunTraceroute             string `json:"runTraceroute" yaml:"runTraceroute"`
	StringToSend              string `json:"stringToSend" yaml:"stringToSend"`
	StringToReceive           string `json:"stringToReceive" yaml:"stringToReceive"`
	Note                      string `json:"note" yaml:"note"`
	UserID                    int    `json:"userId" yaml:"userId"`
	MonitorIntervalPolicy     string `json:"monitorIntervalPolicy" yaml:"monitorIntervalPolicy"`
	NotificationGroups        []int  `json:"notificationGroups" yaml:"notificationGroups"`
	ScheduleID                int    `json:"scheduleId" yaml:"scheduleId"`
	NotificationReportTimeout int    `json:"notificationReportTimeout" yaml:"notificationReportTimeout"`
	VerificationPolicy        string `json:"verificationPolicy" yaml:"verificationPolicy"`
}

type ExpectedSonarTCPCheck struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	SonarTCPCheck
	// This fields must be present when running update action
	mandatoryStructFields []string
}

// UnmarshalYAML unmarshals the mesage and stores original fields
func (ex *ExpectedSonarTCPCheck) UnmarshalYAML(value *yaml.Node) error {

	// Unmarshall data into SonarTCPCheck struct
	var s SonarTCPCheck
	err := value.Decode(&s)
	if err != nil {
		return err
	}
	ex.SonarTCPCheck = s

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
	ex.definedFieldsMap = getFieldNamesMap(&ex.SonarTCPCheck, "yaml", definedFields...)
	ex.mandatoryStructFields = []string{"Name", "ProtocolType", "Port"}

	return nil
}

// Compare compares ExpectedSonarTCPCheck with active SonarTCPCheck. Returns Status and data
func (e *ExpectedSonarTCPCheck) Compare(activeResources *[]SonarTCPCheck) (ResourceAction, []byte, error) {
	var active SonarTCPCheck
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
		action = ActionUpate
		expectedValue := reflect.ValueOf(e.SonarTCPCheck)
		activeValue := reflect.ValueOf(active)
		for _, structFieldName := range e.definedFieldsMap {
			fieldExpected := expectedValue.FieldByName(structFieldName)
			fieldActive := activeValue.FieldByName(structFieldName)
			// Compare field values
			if !reflect.DeepEqual(fieldExpected.Interface(), fieldActive.Interface()) {
				diffStructFields = append(diffStructFields, structFieldName)
			}
		}

		if len(diffStructFields) == 0 {
			return ActionOK, nil, nil
		}
	}

	diffStructFields = append(diffStructFields, e.mandatoryStructFields...)

	diffJSONFields := make([]string, 0)
	for k, v := range e.definedFieldsMap {
		if slices.Contains(diffStructFields, v) {
			diffJSONFields = append(diffJSONFields, k)
		}
	}

	dataBytes, err := toFilteredJSON(e.SonarTCPCheck, diffJSONFields...)
	if err != nil {
		return "", nil, err
	}

	return action, dataBytes, nil
}

func (e *ExpectedSonarTCPCheck) GetActive(activeResources *[]SonarTCPCheck) (*SonarTCPCheck, bool) {
	for _, el := range *activeResources {
		if el.Name == e.Name {
			return &el, true
		}
	}
	return nil, false
}

// GetSonarTCPChecks returns active Sonar Checks
func GetSonarTCPChecks() (*[]SonarTCPCheck, error) {
	// Fetch TCP checks
	fmt.Println("Retrieving Sonar TCP Checks...")
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp")
	if err != nil {
		return nil, err
	}
	data, err := makeAPIRequest("GET", endpoint, nil, 200)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sonar TCP checks: %s", err)
	}

	checks := make([]SonarTCPCheck, 0)
	err = json.Unmarshal(data, &checks)
	if err != nil {
		return nil, err
	}
	return &checks, nil
}

func CreateSonarTCPCheck(payload []byte) error {
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp")
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeAPIRequest("POST", endpoint, payloadReader, 201)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to create Sonar TCP checks: %s", err)
	}
	return nil
}

func UpdateSonarTCPCheck(payload []byte, id int) error {
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp", fmt.Sprint(id))
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeAPIRequest("PUT", endpoint, payloadReader, 200)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to update Sonar TCP checks: %s", err)
	}
	return nil
}

func DeleteSonarTCPCheck(payload []byte, id int) error {
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp", fmt.Sprint(id))
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeAPIRequest("DELETE", endpoint, payloadReader, 200)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to delete Sonar TCP checks: %s", err)
	}
	return nil
}
