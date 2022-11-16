package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MainConfig struct {
	Constellix struct {
		Sonar SonarConfig `yaml:"sonar"`
	} `yaml:"constellix"`
}

type SonarConfig struct {
	HTTPChecksConfigFiles []string `yaml:"http_checks"`
	TCPChecksConfigFiles  []string `yaml:"tcp_checks"`
}

type Config struct {
	SonarHTTPChecks []*ExpectedSonarHTTPCheck
	SonarTCPChecks  []*ExpectedSonarTCPCheck
}

func getConfig(configFile string) (*Config, error) {
	// Read configuration file
	if rootVerbose {
		fmt.Printf("Reading configuration file %s...\n", configFile)
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
			fmt.Printf("  reading %s...\n", configToRead)
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
		}
	}
	for _, item := range mainConfig.Constellix.Sonar.TCPChecksConfigFiles {
		configToRead := filepath.Join(filepath.Dir(configFile), item)
		if rootVerbose {
			fmt.Printf("  reading %s...\n", configToRead)
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
		fmt.Printf("Sonar HTTP Checks saved to %s\n", outputFile)
	} else {
		fmt.Println(string(dataBytes))
	}
	return nil
}
