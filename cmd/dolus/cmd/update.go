/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the init command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Download and install cue definitions for expectation files to use",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	// # curl -s https://raw.githubusercontent.com/DolusExpectation/dolus-expectations/main/install.sh | bash
	//
}
