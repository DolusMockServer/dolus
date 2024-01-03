/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// expectationsCmd represents the expectations command
var expectationsCmd = &cobra.Command{
	Use:   "expectations",
	Short: "Expectations commands",
	Long:  `Expectations commands`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(expectationsCmd)
}
