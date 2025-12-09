package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ppx",
	Short: "PPX - Go web project scaffolding tool",
	Long: `PPX is a CLI tool for quickly creating Go web applications
with clean architecture and production-ready features.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}