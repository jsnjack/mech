/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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

		httpChecks, err := GetSonarChecks()
		if err != nil {
			return err
		}

		fmt.Printf("Found %d Sonar HTTP Checks\n", len(*httpChecks))
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
		inputFile, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		if inputFile == "" {
			return fmt.Errorf("provide configuration file location via --input argument")
		}

		dataBytes, err := os.ReadFile(inputFile)
		if err != nil {
			return err
		}

		fmt.Println(string(dataBytes))
		config := make([]ExpectedSonarHTTPCheck, 0)
		err = yaml.Unmarshal(dataBytes, &config)
		if err != nil {
			return err
		}
		if len(config) == 0 {
			fmt.Println("configuration is empty, nothing to do")
			return nil
		}

		httpChecks, err := GetSonarChecks()
		if err != nil {
			return err
		}

		for _, expectedCheck := range config {
			fmt.Printf("%+v\n", expectedCheck)
			action, _, err := expectedCheck.Compare(httpChecks)
			if err != nil {
				return err
			}
			fmt.Printf("%s: %s\n", action, expectedCheck.Name)
			switch action {
			case ActionOK:
			case ActionDelete:
			case ActionUpate:
			case ActionCreate:
			default:
				return fmt.Errorf("unhandled action %q", action)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sonarCmd)
	sonarCmd.AddCommand(sonarDiscoverCmd)
	sonarDiscoverCmd.PersistentFlags().StringP("output", "o", "", "write output in yaml format to file, filepath")

	sonarCmd.AddCommand(sonarSyncCmd)
	sonarSyncCmd.PersistentFlags().StringP("input", "i", "", "configuration file, filepath")
	sonarSyncCmd.PersistentFlags().Bool("doit", false, "apply planned changes")
	sonarSyncCmd.PersistentFlags().Bool("remove", false, "remove resources which are not present in configuration file")
}
