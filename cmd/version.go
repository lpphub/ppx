package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of PPX",
	Long:  `All software has versions. This is PPX's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PPX v1.0.0")
	},
}