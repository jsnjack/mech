package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

// writeTestFiles creates each file (path relative to dir -> content) under dir.
func writeTestFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestReadConfigs(t *testing.T) {
	dir := t.TempDir()
	writeTestFiles(t, dir, map[string]string{
		"a.yaml":     "a",
		"sub/b.yaml": "b",
		"sub/c.yaml": "c",
	})

	tests := []struct {
		name     string
		baseDir  string
		paths    []string
		expected []string
	}{
		{"literal file", dir, []string{"a.yaml"}, []string{"a"}},
		{"glob with several matches", dir, []string{"sub/*.yaml"}, []string{"b", "c"}},
		{"glob without matches", dir, []string{"none/*.yaml"}, nil},
		{"literal file and glob", dir, []string{"a.yaml", "sub/*.yaml"}, []string{"a", "b", "c"}},
		{"path relative to base dir", filepath.Join(dir, "sub"), []string{"b.yaml"}, []string{"b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := readConfigs(tt.paths, tt.baseDir)
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, item := range data {
				got = append(got, string(item))
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	validFiles := map[string]string{
		"config.yaml": `
constellix:
  sonar:
    http_checks:
      - sonar/http.yaml
    tcp_checks:
      - sonar/tcp/*.yaml
  geoproximity:
    - geoproximity.yaml
  dns:
    example.com:
      - domains/example.com/*.yaml
`,
		"sonar/http.yaml": `
- name: web
  host: web.example.com
  ipVersion: IPV4
  port: 443
  protocolType: HTTPS
  interval: ONEMINUTE
  checkSites: [1, 2]
`,
		"sonar/tcp/a.yaml": `
- name: db
  host: db.example.com
  ipVersion: IPV4
  port: 5432
  interval: ONEMINUTE
  checkSites: [1]
`,
		"sonar/tcp/b.yaml": `
- name: cache
  host: cache.example.com
  ipVersion: IPV4
  port: 6379
  interval: ONEMINUTE
  checkSites: [1]
`,
		"geoproximity.yaml": `
- name: Paris
  longitude: 2.35
  latitude: 48.86
`,
		"domains/example.com/www.yaml": `
- name: www
  type: A
  ttl: 60
  mode: standard
  region: default
  enabled: true
  value:
    - value: 1.1.1.1
      enabled: true
`,
	}

	tests := []struct {
		name       string
		configFile string
		overrides  map[string]string
		expectErr  string
		check      func(t *testing.T, c *Config)
	}{
		{
			name:       "all resource types from literal paths and globs",
			configFile: "config.yaml",
			check: func(t *testing.T, c *Config) {
				if len(c.SonarHTTPChecks) != 1 || c.SonarHTTPChecks[0].Name != "web" {
					t.Fatalf("expected http check %q, got %+v", "web", c.SonarHTTPChecks)
				}
				if !reflect.DeepEqual(c.SonarHTTPChecks[0].CheckSites, []int{1, 2}) {
					t.Errorf("expected checkSites %v, got %v", []int{1, 2}, c.SonarHTTPChecks[0].CheckSites)
				}
				var tcpNames []string
				for _, check := range c.SonarTCPChecks {
					tcpNames = append(tcpNames, check.Name)
				}
				if !reflect.DeepEqual(tcpNames, []string{"db", "cache"}) {
					t.Fatalf("expected tcp checks %q, got %q", []string{"db", "cache"}, tcpNames)
				}
				if c.SonarTCPChecks[0].Port != 5432 {
					t.Errorf("expected port %d, got %d", 5432, c.SonarTCPChecks[0].Port)
				}
				if len(c.GeoProximities) != 1 || c.GeoProximities[0].Latitude != 48.86 {
					t.Errorf("expected geoproximity with latitude 48.86, got %+v", c.GeoProximities)
				}
				records := c.DNS["example.com"]
				if len(records) != 1 || records[0].Name != "www" || records[0].TTL != 60 {
					t.Errorf("expected record www with ttl 60, got %+v", records)
				}
			},
		},
		{
			name:       "glob without matches",
			configFile: "config.yaml",
			overrides: map[string]string{
				"config.yaml": `
constellix:
  sonar:
    tcp_checks:
      - sonar/none/*.yaml
`,
			},
			check: func(t *testing.T, c *Config) {
				if len(c.SonarTCPChecks) != 0 {
					t.Errorf("expected no tcp checks, got %d", len(c.SonarTCPChecks))
				}
			},
		},
		{
			name:       "missing mandatory field",
			configFile: "config.yaml",
			overrides: map[string]string{
				"sonar/tcp/b.yaml": `
- name: cache
  host: cache.example.com
  ipVersion: IPV4
  interval: ONEMINUTE
  checkSites: [1]
`,
			},
			expectErr: `cache: mandatory field "port" is not defined`,
		},
		{
			name:       "malformed sub-config",
			configFile: "config.yaml",
			overrides:  map[string]string{"sonar/http.yaml": "- name: ["},
			expectErr:  "parse sonar http check config",
		},
		{
			name:       "malformed top-level config",
			configFile: "config.yaml",
			overrides:  map[string]string{"config.yaml": "constellix: ["},
			expectErr:  "parse configuration file",
		},
		{
			name:       "missing top-level config",
			configFile: "missing.yaml",
			expectErr:  "read configuration file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeTestFiles(t, dir, validFiles)
			writeTestFiles(t, dir, tt.overrides)

			c, err := getConfig(filepath.Join(dir, tt.configFile))
			if tt.expectErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			tt.check(t, c)
		})
	}
}

func TestWriteDiscoveryResult(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name: "sonar tcp checks",
			value: []*SonarTCPCheck{
				{ID: 1, Name: "db", Host: "db.example.com", Port: 5432, CheckSites: []int{1, 2}, NotificationGroups: []int{}},
			},
			expected: `- id: 1
  name: db
  host: db.example.com
  ipVersion: ""
  port: 5432
  interval: ""
  checkSites:
    - 1
    - 2
  runTraceroute: ""
  stringToSend: ""
  stringToReceive: ""
  note: ""
  userId: 0
  monitorIntervalPolicy: ""
  notificationGroups: []
  scheduleId: 0
  notificationReportTimeout: 0
  verificationPolicy: ""
`,
		},
		{
			name: "geoproximities",
			value: []*GeoProximity{
				{ID: 7, Name: "Paris", Longitude: 2.35, Latitude: 48.86},
			},
			expected: `- id: 7
  name: Paris
  country: ""
  region: ""
  city: 0
  longitude: 2.35
  latitude: 48.86
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFile := filepath.Join(t.TempDir(), "out.yaml")
			if err := writeDiscoveryResult(tt.value, outputFile); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expected, data)
			}

			// The written file must parse back into the same value
			got := reflect.New(reflect.TypeOf(tt.value))
			if err := yaml.Unmarshal(data, got.Interface()); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.Elem().Interface(), tt.value) {
				t.Errorf("expected %+v after round trip, got %+v", tt.value, got.Elem().Interface())
			}
		})
	}
}
