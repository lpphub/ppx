package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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

The module will be created in the 'modules/' directory with:
  - module.go     - Module initialization and route registration
  - model.go      - Database model
  - dto.go        - Data transfer objects (request/response)
  - handler.go    - HTTP handlers
  - service.go    - Business logic
  - repo.go       - Database operations

Examples:
  ppx module product
  ppx module order

After creating the module, you need to:
  1. Import the module in server/app.go
  2. Initialize the module in registerModules()
  3. Run 'go mod tidy' to download dependencies`,
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

	modulePath := filepath.Join("modules", moduleName)
	if _, err := os.Stat(modulePath); err == nil {
		color.Red("❌ Module '%s' already exists", moduleName)
		color.Yellow("💡 Please choose a different module name")
		os.Exit(1)
	}

	moduleDir := "modules"
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		color.Red("❌ 'modules' directory not found")
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

	printModuleSuccess(moduleName, projectModule)
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

func printModuleSuccess(moduleName, projectModule string) {
	color.Green("\n✓ Module '%s' created successfully!", moduleName)
	color.Cyan("\n📂 Generated files:")
	fmt.Printf("   modules/%s/\n", moduleName)
	fmt.Printf("   ├── module.go\n")
	fmt.Printf("   ├── model.go\n")
	fmt.Printf("   ├── dto.go\n")
	fmt.Printf("   ├── handler.go\n")
	fmt.Printf("   ├── service.go\n")
	fmt.Printf("   └── repo.go\n")

	color.Yellow("\n⚠ Next steps:")
	fmt.Printf("   1. Add import \"%s/modules/%s\" to server/app.go\n", projectModule, moduleName)
	fmt.Printf("   2. Initialize the module in registerModules():\n")
	fmt.Printf("      %sMod := %s.New(infra.DB)\n", moduleName, toCamelCase(moduleName))
	fmt.Printf("   3. Add to the modules slice: %sMod\n", moduleName)
	fmt.Printf("   4. Run: go mod tidy && go run .\n")

	color.Cyan("\n📚 Routes will be available at: /api/%s", moduleName)
}
