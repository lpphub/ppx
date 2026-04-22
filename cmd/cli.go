package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ppx",
	Short: "Go web project scaffolding CLI",
	Long: `ppx is a CLI tool for generating Go web projects with modular architecture.

Features:
  - Modular design: each module implements core.Module interface
  - Clean layered architecture: Handler → Service → Repository
  - Built-in JWT authentication with login/register/refresh
  - Demo CRUD module (post) showing best practices
  - Graceful shutdown, health check, and metrics endpoints

Generated Structure:
  config/       - YAML configuration (DB, Redis, JWT, Server)
  modules/      - Business modules (auth, user, post, core)
  infra/        - Infrastructure (DB connections, JWT utils)
  server/       - HTTP server, middleware, helpers
  shared/       - Shared utilities (errors, pagination, contracts)

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
