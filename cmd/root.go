/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootVerbose bool
var rootConfig string
var constellixAPIKey string
var constellixSecretKey string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mech",
	Short: "Constellix DNS configuration as code",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootVerbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().StringVarP(&rootConfig, "filename", "f", "", "configuration filename")
	constellixAPIKey = os.Getenv("CONSTELLIX_API_KEY")
	constellixSecretKey = os.Getenv("CONSTELLIX_SECRET_KEY")
	if constellixAPIKey == "" || constellixSecretKey == "" {
		fmt.Println("Provide CONSTELLIX_API_KEY and CONSTELLIX_SECRET_KEY environmental variables")
		os.Exit(1)
	}
}
