package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ppx",
	Short: "A CLI tool for Go web project scaffolding",
	Long:  `ppx is a command-line tool that helps you quickly create Go web project structures with predefined directories.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
}
