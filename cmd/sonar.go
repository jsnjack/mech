/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

const sonarBaseURL = "https://api.sonar.constellix.com/rest/api"

// sonarCmd represents the sonar command
var sonarCmd = &cobra.Command{
	Use:   "sonar",
	Short: "Sonar checks",
}

// sonarDiscoverCmd represents the discover sonar command
var sonarDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Retrieve existing Sonar configuration",
	Long: `List of supported records:
  - http check
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// Fetch HTTP checks
		fmt.Println("Retrieving Sonar HTTP Checks...")
		endpoint, err := url.JoinPath(sonarBaseURL, "http")
		if err != nil {
			return err
		}
		data, err := makeAPIRequest("GET", endpoint, nil)
		if err != nil {
			return fmt.Errorf("unable to retrieve Sonar HTTP checks: %s", err)
		}

		httpChecks := make([]SonarHTTPCheck, 0)
		err = json.Unmarshal(data, &httpChecks)
		if err != nil {
			return err
		}

		fmt.Printf("Found %d Sonar HTTP Checks\n", len(httpChecks))
		httpCheckBytes, err := yaml.Marshal(httpChecks)
		if err != nil {
			return err
		}

		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		if outputFile != "" {
			err = os.WriteFile(outputFile, httpCheckBytes, 0644)
			if err != nil {
				return err
			}
			fmt.Printf("Sonar HTTP Checks saved to %s\n", outputFile)
		} else {
			fmt.Println(string(httpCheckBytes))
		}

		return nil
	},
}

// sonarSyncCmd represents the sync sonar command
var sonarSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configuration to Constellix",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sonarCmd)
	sonarCmd.AddCommand(sonarDiscoverCmd)
	sonarDiscoverCmd.PersistentFlags().StringP("output", "o", "", "write output in yaml format to file, filepath")

	sonarCmd.AddCommand(sonarSyncCmd)
}
