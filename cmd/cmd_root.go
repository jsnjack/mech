/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootVerbose bool
var rootDebug bool
var constellixAPIKey string
var constellixSecretKey string

var logLevel int

var logger *log.Logger
var reportToTestBuffer bool
var testBuffer *bytes.Buffer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mech",
	Short: "Constellix DNS configuration as code",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set up log level. >0 means verbose, >1 means debug
		fmt.Println("execute")
		if rootVerbose {
			logLevel += 1
		}

		if rootDebug {
			logLevel += 2
		}
	},
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
	logger = log.New(os.Stdout, "", 0)
	testBuffer = new(bytes.Buffer)
	rootCmd.PersistentFlags().BoolVarP(&rootVerbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&rootDebug, "debug", "d", false, "enable debug logging")
	constellixAPIKey = os.Getenv("CONSTELLIX_API_KEY")
	constellixSecretKey = os.Getenv("CONSTELLIX_SECRET_KEY")
	if constellixAPIKey == "" || constellixSecretKey == "" {
		logger.Println("provide CONSTELLIX_API_KEY and CONSTELLIX_SECRET_KEY environmental variables")
		os.Exit(1)
	}
}
