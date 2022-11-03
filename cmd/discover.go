/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

const sonarBaseURL = "https://api.sonar.constellix.com/rest/api"

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Generate configuration file from existing resources",
}

// discoverSonarCmd represents the discover sonar command
var discoverSonarCmd = &cobra.Command{
	Use:   "sonar",
	Short: "Generate configuration file from existing Sonar checks",
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
	discoverCmd.AddCommand(discoverSonarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
