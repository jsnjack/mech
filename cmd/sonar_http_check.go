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

func (ac *SonarHTTPCheck) GetResource() interface{} {
	return ac
}

func (ac *SonarHTTPCheck) GetUID() string {
	return ac.Name
}

func (ac *SonarHTTPCheck) GetConstellixID() int {
	return ac.ID
}

type ExpectedSonarHTTPCheck struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	// List of immutable fields which can't be updated via API
	immutableFields []string
	SonarHTTPCheck
}

// UnmarshalYAML unmarshals the mesage and stores original fields
func (ex *ExpectedSonarHTTPCheck) UnmarshalYAML(value *yaml.Node) error {
	ex.immutableFields = []string{"Host", "IPVersion"}

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
	return nil
}

// GetDefinedStructFieldNames returns list of defined struct fields from local configuration
func (ex *ExpectedSonarHTTPCheck) GetDefinedStructFieldNames() []string {
	return maps.Values(ex.definedFieldsMap)
}

func (ex *ExpectedSonarHTTPCheck) generateData(immutable ...string) ([]byte, error) {
	objBytes, err := json.Marshal(ex)
	if err != nil {
		return nil, err
	}

	// Convert obj to map to simplify iteration
	dataIn := map[string]interface{}{}
	json.Unmarshal(objBytes, &dataIn)

	dataOut := map[string]interface{}{}

	// Create a new data obj which contains only fields which need to be included
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

func (ex *ExpectedSonarHTTPCheck) GetCreateData() ([]byte, error) {
	return ex.generateData()
}

func (ex *ExpectedSonarHTTPCheck) GetUpdateData() ([]byte, error) {
	return ex.generateData(ex.immutableFields...)
}

func (ex *ExpectedSonarHTTPCheck) GetResource() interface{} {
	return ex.SonarHTTPCheck
}

func (ex *ExpectedSonarHTTPCheck) GetUID() string {
	return ex.Name
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
func GetSonarHTTPChecks() ([]*SonarHTTPCheck, error) {
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

	checks := make([]*SonarHTTPCheck, 0)
	err = json.Unmarshal(data, &checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
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
