package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/modules/*.tmpl
var moduleTemplateFS embed.FS

type ModuleData struct {
	ModuleName string
	StructName string
}

func CreateModule(moduleName, structName string) error {
	modulePath := filepath.Join("module", moduleName)
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	data := ModuleData{
		ModuleName: moduleName,
		StructName: structName,
	}

	files := map[string]string{
		"init.go.tmpl":       "init.go",
		"model.go.tmpl":      "model.go",
		"dto.go.tmpl":        "dto.go",
		"handler.go.tmpl":    "handler.go",
		"service.go.tmpl":    "service.go",
		"repository.go.tmpl": "repository.go",
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
