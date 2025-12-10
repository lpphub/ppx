package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// TemplateData holds the data for template rendering
type TemplateData struct {
	ProjectName string
	ModuleName  string
}

// CreateProject creates a new Go web project with the given name
func CreateProject(projectName string) error {
	// Create project root directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Define directories to create
	directories := []string{
		"infra/config",
		"infra/dbs",
		"infra/jwt",
		"infra/logger",
		"infra/monitor",
		"logic",
		"web/handlers",
		"web/middleware",
		"config",
	}

	// Create all directories
	for _, dir := range directories {
		dirPath := filepath.Join(projectName, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Prepare template data
	templateData := TemplateData{
		ProjectName: projectName,
		ModuleName:  projectName,
	}

	// Process templates
	if err := processTemplates(projectName, templateData); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}

	return nil
}

func processTemplates(projectName string, data TemplateData) error {
	// Process main.go template
	if err := processTemplate(
		"templates/cmd/run.go.tmpl",
		filepath.Join(projectName, "main.go"),
		data,
	); err != nil {
		return fmt.Errorf("failed to process main.go template: %w", err)
	}

	// Process web/app.go template
	if err := processTemplate(
		"templates/web/app.go.tmpl",
		filepath.Join(projectName, "web", "app.go"),
		data,
	); err != nil {
		return fmt.Errorf("failed to process web/app.go template: %w", err)
	}

	// Process infra templates
	infraTemplates := map[string]string{
		"templates/infra/config/config.go.tmpl":   "infra/config/config.go",
		"templates/infra/dbs/db.go.tmpl":          "infra/dbs/db.go",
		"templates/infra/jwt/jwt.go.tmpl":         "infra/jwt/jwt.go",
		"templates/infra/logger/logx.go.tmpl":     "infra/logger/logx.go",
		"templates/infra/monitor/monitor.go.tmpl": "infra/monitor/monitor.go",
		"templates/infra/init.go.tmpl":            "infra/init.go",
	}

	for templatePath, outputPath := range infraTemplates {
		if err := processTemplate(
			templatePath,
			filepath.Join(projectName, outputPath),
			data,
		); err != nil {
			return fmt.Errorf("failed to process template %s: %w", templatePath, err)
		}
	}

	// Process config/config.yml template
	if err := processTemplate(
		"templates/config/config.yml.tmpl",
		filepath.Join(projectName, "config", "config.yml"),
		data,
	); err != nil {
		return fmt.Errorf("failed to process config/config.yml template: %w", err)
	}

	return nil
}

func processTemplate(templatePath, outputPath string, data interface{}) error {
	// Read template file
	templateContent, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}
