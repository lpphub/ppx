package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ppx",
	Short: "A CLI tool for Go web project scaffolding",
	Long: `ppx is a command-line tool that helps you quickly create Go web project 
structures with modular architecture, including user authentication, 
and a demo CRUD module.

Project Structure:
  - Modular architecture with module/, contract/, infra/, server/, shared/
  - Built-in user and auth modules
  - Demo post module with full CRUD operations
  - JWT authentication middleware
  - Clean architecture with repository, service, handler layers`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(moduleCmd)
}
