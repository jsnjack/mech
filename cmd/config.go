package cmd

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MainConfig struct {
	Constellix struct {
		Sonar SonarConfig         `yaml:"sonar"`
		DNS   map[string][]string `yaml:"dns"`
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
	for _, item := range mainConfig.Constellix.Sonar.HTTPChecksConfigFiles {
		configToRead := filepath.Join(filepath.Dir(configFile), item)
		if rootVerbose {
			logger.Printf("  reading %s...\n", configToRead)
		}
		configToReadBytes, err := os.ReadFile(configToRead)
		if err != nil {
			return nil, err
		}
		var httpChecks []*ExpectedSonarHTTPCheck
		err = yaml.Unmarshal(configToReadBytes, &httpChecks)
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
	for _, item := range mainConfig.Constellix.Sonar.TCPChecksConfigFiles {
		configToRead := filepath.Join(filepath.Dir(configFile), item)
		if rootVerbose {
			logger.Printf("  reading %s...\n", configToRead)
		}
		configToReadBytes, err := os.ReadFile(configToRead)
		if err != nil {
			return nil, err
		}
		var tcpChecks []*ExpectedSonarTCPCheck
		err = yaml.Unmarshal(configToReadBytes, &tcpChecks)
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
	for domainName, item := range mainConfig.Constellix.DNS {
		for _, recordsFile := range item {
			configToRead := filepath.Join(filepath.Dir(configFile), recordsFile)
			if rootVerbose {
				logger.Printf("  reading %s...\n", configToRead)
			}
			configToReadBytes, err := os.ReadFile(configToRead)
			if err != nil {
				return nil, err
			}
			var records []*ExpectedDNSRecord
			err = yaml.Unmarshal(configToReadBytes, &records)
			if err != nil {
				return nil, err
			}
			config.DNS[domainName] = append(config.DNS[domainName], records...)
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
