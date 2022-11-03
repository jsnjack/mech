/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Generate configuration file from existing resources",
}

// discoverSonarCmd represents the discover sonar command
var discoverSonarCmd = &cobra.Command{
	Use:   "sonar",
	Short: "Generate configuration file from existing Sonar checks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("discover socar called")
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
