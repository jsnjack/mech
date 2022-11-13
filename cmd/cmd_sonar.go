/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/juju/ansiterm"
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

		// Read configuration file
		dataBytes, err := os.ReadFile(inputFile)
		if err != nil {
			return err
		}

		config := make([]ExpectedSonarHTTPCheck, 0)
		err = yaml.Unmarshal(dataBytes, &config)
		if err != nil {
			return err
		}
		if len(config) == 0 {
			fmt.Println("configuration is empty, nothing to do")
			return nil
		}

		// Retrieve active configuration to compare
		httpChecks, err := GetSonarHTTPChecks()
		if err != nil {
			return err
		}

		report := ansiterm.NewTabWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
		defer report.Flush()

		// Check if anything needs to be created / updated
		for _, expectedCheck := range config {
			fmt.Printf("Inspecting %q...\n", expectedCheck.Name)
			action, data, err := expectedCheck.Compare(httpChecks)
			if err != nil {
				return err
			}
			fmt.Printf("  status: %s\n", action)
			fmt.Fprintf(report, "%s\t%s\t%s\n", colorAction(action), expectedCheck.Name, string(data))
			if doit {
				switch action {
				case ActionOK:
					break
				case ActionUpate:
					fmt.Printf("  updating resource %q\n", expectedCheck.Name)
					active, found := expectedCheck.GetActive(httpChecks)
					if found {
						err = UpdateSonarHTTPCheck(data, active.ID)
						if err != nil {
							return err
						}
					} else {
						return fmt.Errorf("%q not found", expectedCheck.Name)
					}
				case ActionCreate:
					fmt.Printf("  creating new resource %q\n", expectedCheck.Name)
					err = CreateSonarHTTPCheck(data)
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("unhandled action %q", action)
				}
			}
		}

		// Check if anything needs to be deleted
		for _, existingCheck := range *httpChecks {
			fmt.Printf("Inspecting %q...\n", existingCheck.Name)
			for _, configCheck := range config {
				if configCheck.Name == existingCheck.Name {
					continue
				}
			}
			fmt.Printf("  status: %s\n", ActionDelete)
			fmt.Fprintf(report, "%s\t%s\t%s\n", colorAction(ActionDelete), existingCheck.Name, "")
			if doit && allowRemoving {
				fmt.Printf("  removing resource %q\n", existingCheck.Name)
				err = DeleteSonarHTTPCheck(nil, existingCheck.ID)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("  pass --remove flag to remove %q\n", existingCheck.Name)
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
