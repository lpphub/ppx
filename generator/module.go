package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/modules/*.tmpl templates/modules/core/module.go.tmpl
var moduleTemplateFS embed.FS

type ModuleData struct {
	ModuleName    string
	StructName    string
	ProjectModule string
}

func CreateModule(moduleName, structName, projectModule string) error {
	modulePath := filepath.Join("modules", moduleName)
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	data := ModuleData{
		ModuleName:    moduleName,
		StructName:    structName,
		ProjectModule: projectModule,
	}

	files := map[string]string{
		"module.go.tmpl":     "module.go",
		"model.go.tmpl":      "model.go",
		"dto.go.tmpl":        "dto.go",
		"handler.go.tmpl":    "handler.go",
		"service.go.tmpl":    "service.go",
		"repo.go.tmpl":       "repo.go",
	}

	for templateName, outputName := range files {
		templateContent, err := moduleTemplateFS.ReadFile("templates/modules/" + templateName)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", templateName, err)
		}

		tmpl, err := template.New(templateName).Parse(string(templateContent))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", templateName, err)
		}

		outputPath := filepath.Join(modulePath, outputName)
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}

		if err := tmpl.Execute(file, data); err != nil {
			file.Close()
			return fmt.Errorf("failed to execute template %s: %w", templateName, err)
		}
		file.Close()
	}

	return nil
}
