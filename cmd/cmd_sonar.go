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

		httpChecks, err := GetSonarHTTPChecks()
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

		// Collect flags
		inputFile, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		if inputFile == "" {
			return fmt.Errorf("provide configuration file location via --input argument")
		}

		doit, err := cmd.Flags().GetBool("doit")
		if err != nil {
			return err
		}

		allowRemoving, err := cmd.Flags().GetBool("remove")
		if err != nil {
			return err
		}

		config, err := getConfig(inputFile)
		if err != nil {
			return err
		}

		// Handle Sonar HTTP Checks
		httpChecks, err := GetSonarHTTPChecks()
		if err != nil {
			return err
		}
		activeHTTPChecks := make([]ResourceMatcher, len(httpChecks))
		for i := range httpChecks {
			activeHTTPChecks[i] = httpChecks[i]
		}
		expectedHTTPChecks := make([]ResourceMatcher, len(config.SonarHTTPChecks))
		for i := range config.SonarHTTPChecks {
			expectedHTTPChecks[i] = &(config.SonarHTTPChecks)[i]
		}
		err = Sync(expectedHTTPChecks, activeHTTPChecks, doit, allowRemoving)
		if err != nil {
			return err
		}
		if !doit {
			fmt.Println("Apply changes by passing --doit flag")
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
