package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ppx %s\n", Version)
		fmt.Printf("  commit: %s\n", GitCommit)
		fmt.Printf("  built:  %s\n", BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}