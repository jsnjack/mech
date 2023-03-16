package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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

func (ac *SonarHTTPCheck) GetResourceID() string {
	return ac.Name
}

func (ac *SonarHTTPCheck) GetConstellixID() int {
	return ac.ID
}

func (ac *SonarHTTPCheck) SyncResourceDelete(constellixID int) error {
	logger.Printf("  removing resource %q\n", ac.GetResourceID())
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http", fmt.Sprint(constellixID))
	if err != nil {
		return err
	}
	body, err := makeSimpleAPIRequest("DELETE", endpoint, nil, 202)
	if err != nil {
		logger.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to delete Sonar HTTP checks: %s", err)
	}
	return nil
}

type ExpectedSonarHTTPCheck struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	// List of immutable fields which can't be updated via API
	immutableFields []string
	// List of mandatory fields which must be defined, used for validation
	mandatoryFields []string
	SonarHTTPCheck
}

// UnmarshalYAML unmarshals the mesage and stores original fields
func (ex *ExpectedSonarHTTPCheck) UnmarshalYAML(value *yaml.Node) error {
	ex.immutableFields = []string{"host", "ipVersion"}
	ex.mandatoryFields = []string{"name", "host", "ipVersion", "port", "protocolType", "interval", "checkSites"}

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

// Validate performs simple validation of user provided data
func (ex *ExpectedSonarHTTPCheck) Validate() error {
	// Validate that all mandatory fields are present
	for _, f := range ex.mandatoryFields {
		if !slices.Contains(maps.Keys(ex.definedFieldsMap), f) {
			return fmt.Errorf("%s: mandatory field %q is not defined", ex.Name, f)
		}
	}
	return nil
}

// GetDefinedStructFieldNames returns list of defined struct fields from local configuration
func (ex *ExpectedSonarHTTPCheck) GetDefinedStructFieldNames() []string {
	return maps.Values(ex.definedFieldsMap)
}

// GetImmutableStructFields returns list of immutable struct fields
func (ex *ExpectedSonarHTTPCheck) GetImmutableStructFields() []string {
	var imf []string
	for k, v := range ex.definedFieldsMap {
		if slices.Contains(ex.immutableFields, k) {
			imf = append(imf, v)
		}
	}
	return imf
}

func (ex *ExpectedSonarHTTPCheck) GetResource() interface{} {
	return ex.SonarHTTPCheck
}

func (ex *ExpectedSonarHTTPCheck) GetResourceID() string {
	return ex.Name
}

func (ex *ExpectedSonarHTTPCheck) SyncResourceUpdate(constellixID int) error {
	logger.Printf("  updating resource %q\n", ex.GetResourceID())
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http", fmt.Sprint(constellixID))
	if err != nil {
		return err
	}
	payload, err := generatePayload(ex, maps.Keys(ex.definedFieldsMap), ex.immutableFields)
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeSimpleAPIRequest("PUT", endpoint, payloadReader, 200)
	if err != nil {
		logger.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to update Sonar HTTP checks: %s", err)
	}
	return nil
}

func (ex *ExpectedSonarHTTPCheck) SyncResourceCreate() error {
	logger.Printf("  creating new resource %q\n", ex.GetResourceID())
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http")
	if err != nil {
		return err
	}
	payload, err := generatePayload(ex, maps.Keys(ex.definedFieldsMap), nil)
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	body, err := makeSimpleAPIRequest("POST", endpoint, payloadReader, 201)
	if err != nil {
		logger.Println("  unexpected response. Details: " + string(body))
		return fmt.Errorf("unable to create Sonar HTTP checks: %s", err)
	}
	return nil
}

// GetSonarHTTPChecks returns active Sonar Checks
func GetSonarHTTPChecks() ([]*SonarHTTPCheck, error) {
	// Fetch HTTP checks
	logger.Println("Retrieving Sonar HTTP Checks...")
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http")
	if err != nil {
		return nil, err
	}
	data, err := makeSimpleAPIRequest("GET", endpoint, nil, 200)
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

// GetSonarHTTPCheckStatus returns active Sonar Check status using runtime endpoint
func GetSonarHTTPCheckStatus(id int) (ResourceRuntimeStatus, error) {
	// Fetch HTTP checks
	logger.Printf("Retrieving status for Sonar HTTP Check %d...\n", id)
	endpoint, err := url.JoinPath(sonarRESTAPIBaseURL, "http", strconv.Itoa(id), "status")
	if err != nil {
		return "", err
	}
	data, err := makeSimpleAPIRequest("GET", endpoint, nil, 200)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Sonar HTTP check status: %s", err)
	}
	status := RuntimeStatus{Status: "unknown"}
	err = json.Unmarshal(data, &status)
	if err != nil {
		return "", err
	}
	return status.Status, nil
}
