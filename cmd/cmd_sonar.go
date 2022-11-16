/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var supportedSonarStaticResources = []string{"http", "tcp"}

// sonarCmd represents the sonar command
var sonarCmd = &cobra.Command{
	Use:   "sonar",
	Short: "Sonar checks",
}

// sonarDiscoverCmd represents the discover sonar command
var sonarDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Fetch Sonar configuration",
}

// sonarDiscoverStaticCmd handle Sonar static commands
// https://api-docs.constellix.com/#01165ee7-fccb-4c96-9fcd-77f329fe6505
var sonarDiscoverStaticCmd = &cobra.Command{
	Use:   "static",
	Short: "Retrieve static configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		resourceType, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}

		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		switch resourceType {
		case "http":
			httpChecks, err := GetSonarHTTPChecks()
			if err != nil {
				return err
			}
			fmt.Printf("Found %d Sonar HTTP Checks\n", len(httpChecks))
			return writeDiscoveryResult(httpChecks, outputFile)
		case "tcp":
			tcpChecks, err := GetSonarTCPChecks()
			if err != nil {
				return err
			}
			fmt.Printf("Found %d Sonar HTTP Checks\n", len(tcpChecks))
			return writeDiscoveryResult(tcpChecks, outputFile)
		default:
			return fmt.Errorf(
				"unsupported resource type: got %q, want one of %q",
				resourceType,
				supportedSonarStaticResources,
			)
		}
	},
}

// sonarSyncCmd represents the sync sonar command
var sonarSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configuration to Constellix",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// Collect flags
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		if configFile == "" {
			return fmt.Errorf("provide configuration file location via --config argument")
		}

		doit, err := cmd.Flags().GetBool("doit")
		if err != nil {
			return err
		}

		allowRemoving, err := cmd.Flags().GetBool("remove")
		if err != nil {
			return err
		}

		config, err := getConfig(configFile)
		if err != nil {
			return err
		}

		// Handle Sonar HTTP Checks
		httpChecks, err := GetSonarHTTPChecks()
		if err != nil {
			return err
		}
		activeHTTPChecks := toResourceMatcher(httpChecks)
		expectedHTTPChecks := toResourceMatcher(config.SonarHTTPChecks)
		err = Sync(expectedHTTPChecks, activeHTTPChecks, doit, allowRemoving)
		if err != nil {
			return err
		}

		// Handle Sonar TCP Checks
		tcpChecks, err := GetSonarTCPChecks()
		if err != nil {
			return err
		}
		activeTCPChecks := toResourceMatcher(tcpChecks)
		expectedTCPChecks := toResourceMatcher(config.SonarTCPChecks)
		err = Sync(expectedTCPChecks, activeTCPChecks, doit, allowRemoving)
		if err != nil {
			return err
		}
		var message string
		if !doit {
			message += "apply changes by passing --doit flag"
		}
		if !allowRemoving {
			if message != "" {
				message += "; "
			}
			message += "allow removing of resources by passing --remove flag"
		}
		if message == "" {
			message = "done"
		}
		fmt.Println(message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sonarCmd)
	sonarCmd.AddCommand(sonarDiscoverCmd)
	sonarDiscoverCmd.PersistentFlags().StringP("output", "o", "", "write output in yaml format to file, filepath")

	sonarDiscoverCmd.AddCommand(sonarDiscoverStaticCmd)
	sonarDiscoverStaticCmd.PersistentFlags().StringP(
		"type", "t", "http", fmt.Sprintf("specify static resource type, one of %q", supportedSonarStaticResources),
	)

	sonarCmd.AddCommand(sonarSyncCmd)
	sonarSyncCmd.PersistentFlags().StringP("config", "c", "", "configuration file, filepath")
	sonarSyncCmd.PersistentFlags().Bool("doit", false, "apply planned changes")
	sonarSyncCmd.PersistentFlags().Bool("remove", false, "remove resources which are not present in configuration file")
}
