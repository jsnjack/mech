package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var dnsRecordResourceIDTemplate = "%s %q (%s, %d)"

// Missing fields: lastValues, skipLookup, contacts
type DNSRecord struct {
	ID                   int              `json:"id"`
	Name                 string           `json:"name" yaml:"name"`
	Type                 string           `json:"type" yaml:"type"`
	TTL                  int              `json:"ttl" yaml:"ttl"`
	Mode                 string           `json:"mode" yaml:"mode"`
	Region               string           `json:"region" yaml:"region"`
	IPFilter             *DNSIPFilter     `json:"ipfilter" yaml:"ipfilter"`
	IPFilterIPDrop       bool             `json:"ipfilteripDrop" yaml:"ipfilteripDrop"`
	GeoFailover          bool             `json:"geoFailover" yaml:"geoFailover"`
	GeoProximity         *DNSGeoProximity `json:"geoProximity" yaml:"geoProximity"`
	Enabled              bool             `json:"enabled" yaml:"enabled"`
	Value                interface{}      `json:"value" yaml:"value"`
	Notes                string           `json:"notes" yaml:"notes"`
	domainIDInConstellix int
}

func (ac *DNSRecord) UnmarshalJSON(b []byte) error {
	var alias aliasDNSRecord
	err := json.Unmarshal(b, &alias)
	if err != nil {
		return err
	}
	s := DNSRecord(alias)
	err = populateDNSRecordValue(&s)
	if err != nil {
		return err
	}
	*ac = s
	return nil
}

func (ac *DNSRecord) GetResource() interface{} {
	return ac
}

func (ac *DNSRecord) GetResourceID() string {
	if ac.GeoProximity != nil {
		return fmt.Sprintf(dnsRecordResourceIDTemplate, ac.Type, ac.Name, ac.Region, ac.GeoProximity.ID)
	}
	return fmt.Sprintf(dnsRecordResourceIDTemplate, ac.Type, ac.Name, ac.Region, 0)
}

func (ac *DNSRecord) GetConstellixID() int {
	return ac.ID
}

func (ac *DNSRecord) SyncResourceDelete(constellixID int) error {
	logger.Printf("  removing resource %q\n", ac.GetResourceID())
	if ac.domainIDInConstellix == 0 {
		return fmt.Errorf("unable to create DNS record: domain ID is not defined (internal error)")
	}
	endpoint, err := url.JoinPath(
		dnsRESTAPIBaseURL,
		"domains",
		fmt.Sprintf("%d", ac.domainIDInConstellix),
		"records",
		fmt.Sprintf("%d", constellixID),
	)
	if err != nil {
		return err
	}
	data, err := makev4APIRequest("DELETE", endpoint, nil, 204)
	if err != nil {
		var details string
		for _, item := range data {
			details += string(item)
		}
		logger.Println("  unexpected response. Details: " + details)
		return fmt.Errorf("unable to delete DNS record: %s", err)
	}
	return nil
}

type DNSIPFilter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DNSGeoProximity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ExpectedDNSRecord struct {
	// Mapping of defined fields from parsed data to struct Field Names
	definedFieldsMap map[string]string
	// List of immutable fields which can't be updated via API
	immutableFields []string
	// List of mandatory fields which must be defined, used for validation
	mandatoryFields []string
	DNSRecord
}

// UnmarshalYAML unmarshals the mesage and stores original fields
func (ex *ExpectedDNSRecord) UnmarshalYAML(value *yaml.Node) error {
	ex.immutableFields = []string{"type"}
	ex.mandatoryFields = []string{"type", "value"}

	// Unmarshall data into DNSRecord struct
	var s DNSRecord
	err := value.Decode(&s)
	if err != nil {
		return err
	}

	err = populateDNSRecordValue(&s)
	if err != nil {
		return err
	}
	ex.DNSRecord = s

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
	ex.definedFieldsMap = getFieldNamesMap(&ex.DNSRecord, "yaml", definedFields...)
	return nil
}

// GetDefinedStructFieldNames returns list of defined struct fields from local configuration
func (ex *ExpectedDNSRecord) GetDefinedStructFieldNames() []string {
	return maps.Values(ex.definedFieldsMap)
}

// GetImmutableStructFields returns list of immutable struct fields
func (ex *ExpectedDNSRecord) GetImmutableStructFields() []string {
	var imf []string
	for k, v := range ex.definedFieldsMap {
		if slices.Contains(ex.immutableFields, k) {
			imf = append(imf, v)
		}
	}
	return imf
}

func (ex *ExpectedDNSRecord) GetResource() interface{} {
	return ex.DNSRecord
}

func (ex *ExpectedDNSRecord) GetResourceID() string {
	if ex.GeoProximity != nil {
		return fmt.Sprintf(dnsRecordResourceIDTemplate, ex.Type, ex.Name, ex.Region, ex.GeoProximity.ID)
	}
	return fmt.Sprintf(dnsRecordResourceIDTemplate, ex.Type, ex.Name, ex.Region, 0)
}

func (ex *ExpectedDNSRecord) SyncResourceUpdate(constellixID int) error {
	logger.Printf("  updating resource %q\n", ex.GetResourceID())
	if ex.domainIDInConstellix == 0 {
		return fmt.Errorf("unable to create DNS record: domain ID is not defined (internal error)")
	}
	endpoint, err := url.JoinPath(
		dnsRESTAPIBaseURL,
		"domains",
		fmt.Sprintf("%d", ex.domainIDInConstellix),
		"records",
		fmt.Sprintf("%d", constellixID),
	)
	if err != nil {
		return err
	}
	payload, err := generatePayload(ex, maps.Keys(ex.definedFieldsMap), ex.immutableFields)
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	data, err := makev4APIRequest("PUT", endpoint, payloadReader, 200)
	if err != nil {
		var details string
		for _, item := range data {
			details += string(item)
		}
		logger.Println("  unexpected response. Details: " + details)
		return fmt.Errorf("unable to update DNS record: %s", err)
	}
	return nil
}

func (ex *ExpectedDNSRecord) SyncResourceCreate() error {
	logger.Printf("  creating new resource %q\n", ex.GetResourceID())
	if ex.domainIDInConstellix == 0 {
		return fmt.Errorf("unable to create DNS record: domain ID is not defined (internal error)")
	}
	endpoint, err := url.JoinPath(dnsRESTAPIBaseURL, "domains", fmt.Sprintf("%d", ex.domainIDInConstellix), "records")
	if err != nil {
		return err
	}
	payload, err := generatePayload(ex, maps.Keys(ex.definedFieldsMap), nil)
	if err != nil {
		return err
	}
	payloadReader := bytes.NewReader(payload)
	data, err := makev4APIRequest("POST", endpoint, payloadReader, 202)
	if err != nil {
		var details string
		for _, item := range data {
			details += string(item)
		}
		logger.Println("  unexpected response. Details: " + details)
		return fmt.Errorf("unable to create DNS record: %s", err)
	}
	return nil
}

// GetDNSRecords retrieves domain's DNS records
func GetDNSRecords(id int) ([]*DNSRecord, error) {
	logger.Printf("Retrieving DNS records for domain %d...\n", id)
	endpoint, err := url.JoinPath(dnsRESTAPIBaseURL, "domains", fmt.Sprintf("%d", id), "records")
	if err != nil {
		return nil, err
	}
	data, err := makev4APIRequest("GET", endpoint, nil, 200)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve DNS domains list: %s", err)
	}

	var records []*DNSRecord
	for _, item := range data {
		var tmpRecords []*DNSRecord
		err = json.Unmarshal(item, &tmpRecords)
		if err != nil {
			return nil, err
		}
		if len(tmpRecords) > 0 {
			records = append(records, tmpRecords...)
		}
	}

	for _, item := range records {
		item.domainIDInConstellix = id
	}
	return records, nil
}
