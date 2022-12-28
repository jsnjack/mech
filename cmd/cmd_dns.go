/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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

		GetDNSDomains()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
	dnsCmd.AddCommand(dnsDiscoverCmd)
	dnsDiscoverCmd.AddCommand(dnsDiscoverRecordsCmd)
}
