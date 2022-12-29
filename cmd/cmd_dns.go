/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "dns configuration",
}

// dnsDiscoverCmd represents the discover sonar command
var dnsDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "fetch configuration from DNS endpoints (domains, records, geo proximities, etc.)",
}

// dnsDiscoverRecordsCmd fetch existing records for a domain name from Constellix
// https://api.dns.constellix.com/v4/docs#tag/Domain-Records/paths/~1domains~1%7Bdomain_id%7D~1records/get
var dnsDiscoverRecordsCmd = &cobra.Command{
	Use:   "records <domain name>",
	Short: "retrieve DNS records for a domain name",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires a domain name (e.g. example.com)")
		}

		if len(args) != 1 {
			return fmt.Errorf("cannot discover multiple domain names at once")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		domains, err := GetDNSDomains()
		if err != nil {
			return err
		}

		var domainID int

		for _, domain := range domains {
			if domain.Name == args[0] {
				domainID = domain.ID
			}
		}

		if domainID == 0 {
			return fmt.Errorf("domain %s not found", args[0])
		} else {
			if rootVerbose {
				logger.Printf("domain %s found with ID %d", args[0], domainID)
			}
		}
		records, err := GetDNSRecords(domainID)
		if err != nil {
			return err
		}
		logger.Printf("Found %d DNS records\n", len(records))
		return writeDiscoveryResult(records, outputFile)
	},
}

// dnsDiscoverDomainsCmd fetch existing domains from Constellix
var dnsDiscoverDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "retrieve domains registered in Constellix",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		domains, err := GetDNSDomains()
		if err != nil {
			return err
		}

		report := table.NewWriter()
		// For tests, render data in csv format
		defer func() {
			if report.Length() > 0 {
				if reportToTestBuffer {
					report.RenderCSV()
				} else {
					report.Render()
				}
			} else {
				logger.Println("  nothing to do")
			}
		}()

		if reportToTestBuffer {
			// Skip header in tests
			report.SetOutputMirror(testBuffer)
		} else {
			report.SetOutputMirror(os.Stdout)
			report.AppendHeader(table.Row{"ID", "Name", "Status", "GeoIP", "GTD"})
		}

		for _, domain := range domains {
			report.AppendRow(table.Row{
				domain.ID, domain.Name, domain.Status, domain.GeoIPEnabled, domain.GTDEnabled,
			})
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
	dnsCmd.AddCommand(dnsDiscoverCmd)
	dnsDiscoverCmd.AddCommand(dnsDiscoverRecordsCmd)
	dnsDiscoverCmd.PersistentFlags().StringP("output", "o", "", "write output in yaml format to file, filepath")

	dnsDiscoverCmd.AddCommand(dnsDiscoverDomainsCmd)
}
