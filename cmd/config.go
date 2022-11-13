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
}

type Config struct {
	SonarHTTPChecks []ExpectedSonarHTTPCheck
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
		var httpChecks []ExpectedSonarHTTPCheck
		err = yaml.Unmarshal(configToReadBytes, &httpChecks)
		if err != nil {
			return nil, err
		}
		if len(httpChecks) > 0 {
			config.SonarHTTPChecks = append(config.SonarHTTPChecks, httpChecks...)
		}
	}
	return &config, nil
}
