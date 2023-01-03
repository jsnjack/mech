package cmd

import (
	"encoding/json"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestExpectedDNSRecord_UnmarshalYAML(t *testing.T) {
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

func TestExpectedDNSRecord_UnmarshalJSON(t *testing.T) {
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
