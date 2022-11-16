package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

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

func (ac *SonarTCPCheck) GetResource() interface{} {
	return ac
}

func (ac *SonarTCPCheck) GetResourceID() string {
	return ac.Name
}

func (ac *SonarTCPCheck) GetConstellixID() int {
	return ac.ID
}

func (ac *SonarTCPCheck) SyncResourceDelete(constellixID int) error {
	fmt.Printf("  removing resource %q\n", ac.GetResourceID())
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp", fmt.Sprint(constellixID))
	if err != nil {
		return err
	}
	body, err := makeAPIRequest("DELETE", endpoint, nil, 202)
	if err != nil {
		fmt.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to delete Sonar TCP checks: %s", err)
	}
	return nil
}

type ExpectedSonarTCPCheck struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	// List of immutable fields which can't be updated via API
	immutableFields []string
	SonarTCPCheck
}

// UnmarshalYAML unmarshals the mesage and stores original fields
func (ex *ExpectedSonarTCPCheck) UnmarshalYAML(value *yaml.Node) error {
	ex.immutableFields = []string{"host", "ipVersion"}

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
	return nil
}

// GetDefinedStructFieldNames returns list of defined struct fields from local configuration
func (ex *ExpectedSonarTCPCheck) GetDefinedStructFieldNames() []string {
	return maps.Values(ex.definedFieldsMap)
}

func (ex *ExpectedSonarTCPCheck) generateData(immutable ...string) ([]byte, error) {
	objBytes, err := json.Marshal(ex)
	if err != nil {
		return nil, err
	}

	// Convert obj to map to simplify iteration
	dataIn := map[string]interface{}{}
	json.Unmarshal(objBytes, &dataIn)

	dataOut := map[string]interface{}{}

	// Create a new data obj which contains only fields which need to be included (JSON)
	for key, value := range dataIn {
		if slices.Contains(maps.Keys(ex.definedFieldsMap), key) && !slices.Contains(immutable, key) {
			dataOut[key] = value
		}
	}

	dataOutBytes, err := json.Marshal(dataOut)
	if err != nil {
		return nil, err
	}
	return dataOutBytes, nil
}

func (ex *ExpectedSonarTCPCheck) GetResource() interface{} {
	return ex.SonarTCPCheck
}

func (ex *ExpectedSonarTCPCheck) GetResourceID() string {
	return ex.Name
}

func (ex *ExpectedSonarTCPCheck) SyncResourceUpdate(constellixID int) error {
	fmt.Printf("  updating resource %q\n", ex.GetResourceID())
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp", fmt.Sprint(constellixID))
	if err != nil {
		return err
	}
	payload, err := ex.generateData(ex.immutableFields...)
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

func (ex *ExpectedSonarTCPCheck) SyncResourceCreate() error {
	fmt.Printf("  creating new resource %q\n", ex.GetResourceID())
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "tcp")
	if err != nil {
		return err
	}
	payload, err := ex.generateData()
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

// GetSonarTCPChecks returns active Sonar Checks
func GetSonarTCPChecks() ([]*SonarTCPCheck, error) {
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

	checks := make([]*SonarTCPCheck, 0)
	err = json.Unmarshal(data, &checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}
