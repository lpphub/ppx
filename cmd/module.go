package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"ppx/generator"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:   "module [module-name]",
	Short: "Create a new module in an existing project",
	Long: `Create a new module with full CRUD structure in an existing project.

The module will be created across MVC directories:
  - handler/<name>.go        - HTTP handlers + DTO
  - repository/models/<name>.go - GORM model
  - repository/<name>.go     - Database operations
  - service/<name>.go        - Business logic

Examples:
  ppx module product
  ppx module order

After creating the module, you need to:
  1. Add handler initialization in route/route.go
  2. Run 'go mod tidy' to download dependencies`,
	Args: cobra.ExactArgs(1),
	Run:  runModule,
}

func runModule(cmd *cobra.Command, args []string) {
	moduleName := args[0]

	if err := validateModuleName(moduleName); err != nil {
		color.Red("❌ Invalid module name: %v", err)
		color.Yellow("💡 Module name must be lowercase letters, numbers, and underscores")
		os.Exit(1)
	}

	// Check if handler directory exists (indicates ppx-generated project)
	if _, err := os.Stat("handler"); os.IsNotExist(err) {
		color.Red("❌ 'handler' directory not found")
		color.Yellow("💡 Please run this command in a ppx-generated project root directory")
		os.Exit(1)
	}

	projectModule := getProjectModule()
	if projectModule == "" {
		color.Red("❌ Could not read project module from go.mod")
		color.Yellow("💡 Make sure go.mod exists in the current directory")
		os.Exit(1)
	}

	moduleName = strings.ToLower(moduleName)
	structName := toCamelCase(moduleName)

	if err := generator.CreateModule(moduleName, structName, projectModule); err != nil {
		color.Red("❌ Failed to create module '%s'", moduleName)
		color.Yellow("💡 Error details: %v", err)
		os.Exit(1)
	}

	printModuleSuccess(moduleName, structName, projectModule)
}

func getProjectModule() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func validateModuleName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("module name cannot be empty")
	}
	if len(name) > 30 {
		return fmt.Errorf("module name too long (max 30 characters)")
	}
	if !regexp.MustCompile(`^[a-z][a-z0-9_]*$`).MatchString(name) {
		return fmt.Errorf("module name must start with a lowercase letter and contain only lowercase letters, numbers, and underscores")
	}
	return nil
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func printModuleSuccess(moduleName, structName, projectModule string) {
	color.Green("\n✓ Module '%s' created successfully!", moduleName)
	color.Cyan("\n📂 Generated files:")
	fmt.Printf("   handler/%s.go\n", moduleName)
	fmt.Printf("   types/models/%s.go\n", moduleName)
	fmt.Printf("   types/%s.go\n", moduleName)
	fmt.Printf("   repository/%s.go\n", moduleName)
	fmt.Printf("   service/%s.go\n", moduleName)

	color.Yellow("\n⚠ Next steps:")
	fmt.Printf("   1. Add to route/route.go:\n")
	fmt.Printf("      %sRepo := repository.New%sRepo(infra.DB)\n", moduleName, structName)
	fmt.Printf("      %sSvc := service.New%sService(%sRepo)\n", moduleName, structName, moduleName)
	fmt.Printf("      %sHandler := handler.New%sHandler(%sSvc)\n", moduleName, structName, moduleName)
	fmt.Printf("      %sHandler.RegisterRoutes(group)\n", moduleName)
	fmt.Printf("   2. Run: go mod tidy && go run .\n")

	color.Cyan("\n📚 Routes will be available at: /%ss", moduleName)
}
