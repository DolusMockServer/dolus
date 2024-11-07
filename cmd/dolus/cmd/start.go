/*
Copyright © 2023 Martin Simango shukomango@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/DolusMockServer/dolus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a mock server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		spec := cmd.Flag("spec").Value.String()
		cueExpectationsFiles, _ := cmd.Flags().GetStringArray("cueExpectationsFiles")
		port, _ := cmd.Flags().GetInt("port")

		d := dolus.New()

		d.AddExpectations(cueExpectationsFiles...)

		d.OpenAPIspec = spec
		if err := d.Start(fmt.Sprintf(":%d", port)); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("spec", "s", "", "openapi spec file")
	startCmd.Flags().StringArrayP("cueExpectationsFiles", "e", []string{}, "cue expectation files")
	startCmd.Flags().IntP("port", "p", 1080, "port to start server on")
	startCmd.MarkFlagRequired("spec")
}
