package cmd

import "testing"

func TestPopulateDNSRecordIpfilterForJSON_ValidObject(t *testing.T) {
	record := &DNSRecord{
		IPFilter: map[string]interface{}{
			"id":   1,
			"name": "test",
		},
	}

	err := populateDNSRecordIPFilterForJSON(record)
	if err != nil {
		t.Errorf("populateDNSRecordIPFilterForJSON() error = %v, wantErr %v", err, false)
	}
	if record.IPFilter != 1 {
		t.Errorf("populateDNSRecordIPFilterForJSON() = %v, want %v", record.IPFilter, 1)
	}
}
