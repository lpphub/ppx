package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ppx",
	Short: "Go web project scaffolding CLI",
	Long: `ppx is a CLI tool for generating Go web projects with MVC architecture.

Features:
  - MVC architecture: handler → service → repository
  - Centralized route registration in routes/
  - Built-in JWT authentication with login/register/refresh
  - Demo CRUD module (post) showing best practices
  - Graceful shutdown, health check, and metrics endpoints

Generated Structure:
  handler/      - HTTP handlers
  service/      - Business logic
  repository/   - Data access
  types/        - Models & DTOs
  routes/       - Route registration
  middleware/   - HTTP middleware
  shared/       - Common utilities
  infra/        - Infrastructure (DB, Redis, JWT)

Commands:
  ppx new [name]      Create a new project
  ppx module [name]   Add a new CRUD module to existing project
  ppx version         Show version info`,
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
