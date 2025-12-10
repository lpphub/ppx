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
		"infra/logger/logx",
		"infra/monitor",
		"logic",
		"logic/auth",
		"logic/user",
		"logic/shared",
		"web/base",
		"web/handlers",
		"web/middleware",
		"web/rest",
		"web/rest/permission",
		"web/rest/role",
		"web/rest/user",
		"web/types",
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
	// All templates mapping
	templates := map[string]string{
		"templates/main.go.tmpl":                    "main.go",
		"templates/web/app.go.tmpl":                 "web/app.go",
		"templates/web/base/render.go.tmpl":         "web/base/render.go",
		"templates/web/middleware/auth.go.tmpl":     "web/middleware/auth.go",
		"templates/web/rest/handler.go.tmpl":        "web/rest/handler.go",
		"templates/web/types/req.go.tmpl":           "web/types/req.go",
		"templates/web/types/resp.go.tmpl":          "web/types/resp.go",
		"templates/infra/config/config.go.tmpl":     "infra/config/config.go",
		"templates/infra/dbs/db.go.tmpl":            "infra/dbs/db.go",
		"templates/infra/jwt/jwt.go.tmpl":           "infra/jwt/jwt.go",
		"templates/infra/logger/logger.go.tmpl":     "infra/logger/logger.go",
		"templates/infra/logger/global.go.tmpl":     "infra/logger/global.go",
		"templates/infra/logger/extractor.go.tmpl":  "infra/logger/extractor.go",
		"templates/infra/logger/zerolog.go.tmpl":    "infra/logger/zerolog.go",
		"templates/infra/logger/logx/gin.go.tmpl":   "infra/logger/logx/gin.go",
		"templates/infra/logger/logx/gorm.go.tmpl":  "infra/logger/logx/gorm.go",
		"templates/infra/logger/logx/redis.go.tmpl": "infra/logger/logx/redis.go",
		"templates/infra/monitor/metrics.go.tmpl":   "infra/monitor/metrics.go",
		"templates/infra/monitor/pprof.go.tmpl":     "infra/monitor/pprof.go",
		"templates/infra/init.go.tmpl":              "infra/init.go",
		"templates/logic/auth/service.go.tmpl":      "logic/auth/service.go",
		"templates/logic/user/service.go.tmpl":      "logic/user/service.go",
		"templates/logic/user/model.go.tmpl":        "logic/user/model.go",
		"templates/logic/shared/consts.go.tmpl":     "logic/shared/consts.go",
		"templates/logic/shared/errors.go.tmpl":     "logic/shared/errors.go",
		"templates/logic/shared/pagination.go.tmpl": "logic/shared/pagination.go",
		"templates/config/config.yml.tmpl":          "config/config.yml",
		"templates/go.mod.tmpl":                     "go.mod",
		"templates/Dockerfile.tmpl":                 "Dockerfile",
	}

	// Process all templates
	for templatePath, outputPath := range templates {
		if err := processTemplate(
			templatePath,
			filepath.Join(projectName, outputPath),
			data,
		); err != nil {
			return fmt.Errorf("failed to process template %s: %w", templatePath, err)
		}
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
