package cmd

import (
	"encoding/json"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestExpectedDNSRecord_Standard_UnmarshalYAML(t *testing.T) {
	data := `
name: abc
type: A
ttl: 60
mode: standard
region: default
enabled: true
value:
  - value: 1.1.1.1
    enabled: true
`
	var obj ExpectedDNSRecord
	err := yaml.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.Name != "abc" {
		t.Errorf("expected %q, got %q", "abc", obj.Name)
		return
	}
	if obj.Mode != "standard" {
		t.Errorf("expected %q, got %q", "standard", obj.Mode)
		return
	}
	expected := DNSStandardValue{
		Value:   "1.1.1.1",
		Enabled: true,
	}
	res, ok := obj.Value.([]*DNSStandardValue)
	if !ok {
		t.Errorf("unexpected type %T", obj.Value)
		return
	}
	if len(res) != 1 {
		t.Errorf("expected 1, got %d", len(res))
		return
	}
	value := res[0]
	if value.Value != expected.Value {
		t.Errorf("expected %q, got %q", expected.Value, value.Value)
		return
	}
	if value.Enabled != expected.Enabled {
		t.Errorf("expected %t, got %t", expected.Enabled, value.Enabled)
		return
	}
}

func TestExpectedDNSRecord_Standard_UnmarshalJSON(t *testing.T) {
	data := `{"id":31847357,"name":"abc","type":"A","ttl":600,"mode":"standard","region":"default","ipfilter":null,"ipfilterDrop":false,"geoFailover":false,"geoproximity":null,"enabled":true,"value":[{"value":"1.1.1.1","enabled":true}],"lastValues":{"roundRobinFailover":[],"standard":[{"value":"8.8.8.8","enabled":true}],"failover":{"enabled":false,"mode":"normal","values":[]},"pools":[]},"notes":"","skipLookup":null,"domain":{"id":1004580,"name":"surfly.gratis","status":"ACTIVE","geoip":false,"gtd":false,"tags":[],"createdAt":"2022-12-28T15:13:57+00:00","updatedAt":"2022-12-29T19:35:29+00:00","links":{"self":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580","records":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580\/records","history":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580\/history","nameservers":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580\/nameservers","analytics":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580\/analytics"}},"contacts":[],"links":{"self":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580\/records\/31847357","domain":"http:\/\/api.dns.constellix.com\/v4\/domains\/1004580"}}`
	var obj DNSRecord
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.Name != "abc" {
		t.Errorf("expected %q, got %q", "abc", obj.Name)
		return
	}
	if obj.Mode != "standard" {
		t.Errorf("expected %q, got %q", "standard", obj.Mode)
		return
	}
	expected := DNSStandardValue{
		Value:   "1.1.1.1",
		Enabled: true,
	}
	res, ok := obj.Value.([]*DNSStandardValue)
	if !ok {
		t.Errorf("unexpected type %T", obj.Value)
		return
	}
	if len(res) != 1 {
		t.Errorf("expected 1, got %d", len(res))
		return
	}
	value := res[0]
	if value.Value != expected.Value {
		t.Errorf("expected %q, got %q", expected.Value, value.Value)
		return
	}
	if value.Enabled != expected.Enabled {
		t.Errorf("expected %t, got %t", expected.Enabled, value.Enabled)
		return
	}
}

func TestExpectedDNSRecord_Failover_UnmarshalYAML(t *testing.T) {
	data := `
name: abc
type: A
ttl: 60
mode: failover
region: default
enabled: true
value:
  enabled: true
  mode: normal
  values:
    - value: 1.1.1.1
      enabled: true
      order: 1
      sonarCheckId: 123
      active: false
      failed: true
      status: DOWN
    - value: 1.1.1.2
      enabled: true
      order: 2
      sonarCheckId: null
      active: false
      failed: true
      status: N/A
`
	var obj ExpectedDNSRecord
	err := yaml.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.Name != "abc" {
		t.Errorf("expected %q, got %q", "abc", obj.Name)
		return
	}
	if obj.Mode != "failover" {
		t.Errorf("expected %q, got %q", "failover", obj.Mode)
		return
	}
	expected := DNSFailoverValue{
		Enabled: true,
		Mode:    "normal",
		Values: []*DNSFailoverItemValue{
			{
				Value:        "1.1.1.1",
				Enabled:      true,
				Order:        1,
				SonarCheckID: 123,
			},
			{
				Value:        "1.1.1.2",
				Enabled:      true,
				Order:        2,
				SonarCheckID: 0,
			},
		},
	}
	res, ok := obj.Value.(*DNSFailoverValue)
	if !ok {
		t.Errorf("unexpected type %T", obj.Value)
		return
	}
	if res.Enabled != expected.Enabled {
		t.Errorf("expected %t, got %t", expected.Enabled, res.Enabled)
		return
	}
	if res.Mode != expected.Mode {
		t.Errorf("expected %q, got %q", expected.Mode, res.Mode)
		return
	}
	if len(res.Values) != 2 {
		t.Errorf("expected 2, got %d", len(res.Values))
		return
	}
	res1 := res.Values[0]
	if res1.Value != expected.Values[0].Value {
		t.Errorf("expected %q, got %q", expected.Values[0].Value, res1.Value)
		return
	}
	if res1.Enabled != expected.Values[0].Enabled {
		t.Errorf("expected %t, got %t", expected.Values[0].Enabled, res1.Enabled)
		return
	}
	if res1.Order != expected.Values[0].Order {
		t.Errorf("expected %d, got %d", expected.Values[0].Order, res1.Order)
		return
	}
	if res1.SonarCheckID != expected.Values[0].SonarCheckID {
		t.Errorf("expected %d, got %d", expected.Values[0].SonarCheckID, res1.SonarCheckID)
		return
	}
	res2 := res.Values[1]
	if res2.SonarCheckID != expected.Values[1].SonarCheckID {
		t.Errorf("expected %d, got %d", expected.Values[1].SonarCheckID, res2.SonarCheckID)
		return
	}
}

func TestExpectedDNSRecord_Failover_UnmarshalJSON(t *testing.T) {
	data := `{"id":31847262,"name":"abc","type":"A","ttl":60,"mode":"failover","region":"default","ipfilter":null,"ipfilterDrop":false,"geoFailover":false,"geoproximity":null,"enabled":true,"value":{"enabled":true,"mode":"normal","values":[{"value":"159.69.18.28","order":1,"sonarCheckId":84874,"enabled":true,"active":false,"failed":true,"status":"DOWN"},{"value":"1.1.1.1","order":2,"sonarCheckId":null,"enabled":true,"active":false,"failed":false,"status":"N\/A"}]}}`
	var obj ExpectedDNSRecord
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.Name != "abc" {
		t.Errorf("expected %q, got %q", "abc", obj.Name)
		return
	}
	if obj.Mode != "failover" {
		t.Errorf("expected %q, got %q", "failover", obj.Mode)
		return
	}
	expected := DNSFailoverValue{
		Enabled: true,
		Mode:    "normal",
		Values: []*DNSFailoverItemValue{
			{
				Value:        "159.69.18.28",
				Enabled:      true,
				Order:        1,
				SonarCheckID: 84874,
			},
			{
				Value:        "1.1.1.2",
				Enabled:      true,
				Order:        2,
				SonarCheckID: 0,
			},
		},
	}
	res, ok := obj.Value.(*DNSFailoverValue)
	if !ok {
		t.Errorf("unexpected type %T", obj.Value)
		return
	}
	if res.Enabled != expected.Enabled {
		t.Errorf("expected %t, got %t", expected.Enabled, res.Enabled)
		return
	}
	if res.Mode != expected.Mode {
		t.Errorf("expected %q, got %q", expected.Mode, res.Mode)
		return
	}
	if len(res.Values) != 2 {
		t.Errorf("expected 2, got %d", len(res.Values))
		return
	}
	res1 := res.Values[0]
	if res1.Value != expected.Values[0].Value {
		t.Errorf("expected %q, got %q", expected.Values[0].Value, res1.Value)
		return
	}
	if res1.Enabled != expected.Values[0].Enabled {
		t.Errorf("expected %t, got %t", expected.Values[0].Enabled, res1.Enabled)
		return
	}
	if res1.Order != expected.Values[0].Order {
		t.Errorf("expected %d, got %d", expected.Values[0].Order, res1.Order)
		return
	}
	if res1.SonarCheckID != expected.Values[0].SonarCheckID {
		t.Errorf("expected %d, got %d", expected.Values[0].SonarCheckID, res1.SonarCheckID)
		return
	}
	res2 := res.Values[1]
	if res2.SonarCheckID != expected.Values[1].SonarCheckID {
		t.Errorf("expected %d, got %d", expected.Values[1].SonarCheckID, res2.SonarCheckID)
		return
	}
}
