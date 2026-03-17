package cmd

import (
	"fmt"
	"os"
	"regexp"

	"ppx/generator"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go web project",
	Long: `Create a new Go web project with modular architecture.

Examples:
  ppx new myapp
  ppx new myapp --module github.com/user/myapp

Generated Project Structure:
  myapp/
  ├── config/
  │   └── config.yml
  ├── infra/
  │   ├── init.go
  │   ├── config.go
  │   ├── database.go
  │   └── jwt/
  │       └── jwt.go
  ├── module/
  │   ├── contract/      # Contract/Interface definitions
  │   ├── auth/          # Authentication module
  │   ├── user/          # User module
  │   └── post/          # Demo post module (CRUD example)
  ├── server/
  │   ├── app.go
  │   ├── helper/
  │   └── middleware/
  ├── shared/
  │   ├── consts/
  │   ├── errs/
  │   ├── mod/
  │   └── strutils/
  ├── main.go
  ├── go.mod
  ├── Makefile
  └── Dockerfile`,
	Args: cobra.ExactArgs(1),
	Run:  runNew,
}

func runNew(cmd *cobra.Command, args []string) {
	projectName := args[0]

	if err := validateProjectName(projectName); err != nil {
		color.Red("❌ Invalid project name: %v", err)
		color.Yellow("💡 Project name must start with a letter and contain only letters, numbers, hyphens, and underscores")
		os.Exit(1)
	}

	if _, err := os.Stat(projectName); err == nil {
		color.Red("❌ Directory '%s' already exists", projectName)
		color.Yellow("💡 Please choose a different name or remove the existing directory")
		os.Exit(1)
	}

	moduleName, _ := cmd.Flags().GetString("module")
	if moduleName == "" {
		moduleName = projectName
	}

	if err := generator.CreateProject(projectName, moduleName); err != nil {
		handleCreateError(err, projectName)
		os.Exit(1)
	}
}

func handleCreateError(err error, projectName string) {
	color.Red("❌ Failed to create project '%s'", projectName)
	color.Yellow("💡 Error details: %v", err)
	color.Cyan("📞 Need help? Visit: https://github.com/lpphub/ppx/issues")
}

func validateProjectName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("project name cannot be empty")
	}
	if len(name) > 50 {
		return fmt.Errorf("project name too long (max 50 characters)")
	}
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]*$`).MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}
	return nil
}

func init() {
	newCmd.Flags().String("module", "", "Module name for the project (e.g., github.com/user/project)")
}
