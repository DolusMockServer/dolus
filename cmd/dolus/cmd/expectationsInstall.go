/*
Copyright © 2024 Martin Simango shukomango@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// expectationsInstallCmd represents the expectationsDownload command
var expectationsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and install cue definitions for expectation files to use",
	Long: `Download and install cue definitions for expectation files to use`,
	Run: func(cmd *cobra.Command, args []string) {

		c := exec.Command("bash", "-c", "curl -s https://raw.githubusercontent.com/DolusMockServer/Dolus/main/install.sh | bash")

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	},
}

func init() {

	expectationsCmd.AddCommand(expectationsInstallCmd)

}
