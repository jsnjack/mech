package cmd

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MainConfig struct {
	Constellix struct {
		Sonar                   SonarConfig         `yaml:"sonar"`
		GeoProximityConfigFiles []string            `yaml:"geoproximity"`
		DNS                     map[string][]string `yaml:"dns"`
	} `yaml:"constellix"`
}

type SonarConfig struct {
	HTTPChecksConfigFiles []string `yaml:"http_checks"`
	TCPChecksConfigFiles  []string `yaml:"tcp_checks"`
}

type Config struct {
	SonarHTTPChecks []*ExpectedSonarHTTPCheck
	SonarTCPChecks  []*ExpectedSonarTCPCheck
	DNS             map[string][]*ExpectedDNSRecord
	GeoProximities  []*ExpectedGeoProximity
}

func getConfig(configFile string) (*Config, error) {
	// Read configuration file
	if rootVerbose {
		logger.Printf("Reading configuration file %s...\n", configFile)
	}
	dataBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var mainConfig MainConfig
	err = yaml.Unmarshal(dataBytes, &mainConfig)
	if err != nil {
		return nil, err
	}

	var config Config
	dataB, err := readConfigs(mainConfig.Constellix.Sonar.HTTPChecksConfigFiles, filepath.Dir(configFile))
	if err != nil {
		return nil, err
	}
	for _, item := range dataB {
		var httpChecks []*ExpectedSonarHTTPCheck
		err = yaml.Unmarshal(item, &httpChecks)
		if err != nil {
			return nil, err
		}
		if len(httpChecks) > 0 {
			config.SonarHTTPChecks = append(config.SonarHTTPChecks, httpChecks...)
			for _, check := range httpChecks {
				err = check.Validate()
				if err != nil {
					return nil, err
				}
			}
		}
	}

	dataB, err = readConfigs(mainConfig.Constellix.Sonar.TCPChecksConfigFiles, filepath.Dir(configFile))
	if err != nil {
		return nil, err
	}
	for _, item := range dataB {
		var tcpChecks []*ExpectedSonarTCPCheck
		err = yaml.Unmarshal(item, &tcpChecks)
		if err != nil {
			return nil, err
		}
		if len(tcpChecks) > 0 {
			config.SonarTCPChecks = append(config.SonarTCPChecks, tcpChecks...)
			for _, check := range tcpChecks {
				err = check.Validate()
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// DNS
	config.DNS = make(map[string][]*ExpectedDNSRecord)
	for domainName, cfs := range mainConfig.Constellix.DNS {
		dataB, err = readConfigs(cfs, filepath.Dir(configFile))
		if err != nil {
			return nil, err
		}
		for _, dataItem := range dataB {
			var records []*ExpectedDNSRecord
			err = yaml.Unmarshal(dataItem, &records)
			if err != nil {
				return nil, err
			}
			config.DNS[domainName] = append(config.DNS[domainName], records...)
		}
	}

	// GeoProximities
	dataB, err = readConfigs(mainConfig.Constellix.GeoProximityConfigFiles, filepath.Dir(configFile))
	if err != nil {
		return nil, err
	}
	for _, item := range dataB {
		var geops []*ExpectedGeoProximity
		err = yaml.Unmarshal(item, &geops)
		if err != nil {
			return nil, err
		}
		if len(geops) > 0 {
			config.GeoProximities = append(config.GeoProximities, geops...)
			for _, check := range geops {
				err = check.Validate()
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &config, nil
}

func writeDiscoveryResult(collection interface{}, outputFile string) error {
	dataBytes, err := yaml.Marshal(collection)
	if err != nil {
		return err
	}
	if outputFile != "" {
		err = os.WriteFile(outputFile, dataBytes, 0644)
		if err != nil {
			return err
		}
		logger.Printf("Discovered items saved to %s\n", outputFile)
	} else {
		logger.Println(string(dataBytes))
	}
	return nil
}

// readConfigs reads all configuration files. If file doesn't exist, assumes it is
// a glob pattern and reads all files matching the pattern.
func readConfigs(configFiles []string, baseDir string) ([][]byte, error) {
	var dataBytes [][]byte
	for _, configFile := range configFiles {
		configToRead := filepath.Join(baseDir, configFile)
		if _, err := os.Stat(configToRead); err == nil {
			// File exists
			if rootVerbose {
				logger.Printf("  reading %s...\n", configToRead)
			}
			data, err := os.ReadFile(configToRead)
			if err != nil {
				return nil, err
			}
			dataBytes = append(dataBytes, data)
		} else {
			// File doesn't exist, assume it is a glob pattern
			if rootVerbose {
				logger.Printf("  assuming %s is a pattern...\n", configToRead)
			}
			files, err := filepath.Glob(configToRead)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				if rootVerbose {
					logger.Printf("  reading %s...\n", file)
				}
				data, err := os.ReadFile(file)
				if err != nil {
					return nil, err
				}
				dataBytes = append(dataBytes, data)
			}
		}
	}
	return dataBytes, nil
}
