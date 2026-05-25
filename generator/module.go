package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/module-tmpl/*.tmpl
var moduleTemplateFS embed.FS

type ModuleData struct {
	ModuleName    string
	StructName    string
	ProjectModule string
}

func CreateModule(moduleName, structName, projectModule string) error {
	data := ModuleData{
		ModuleName:    moduleName,
		StructName:    structName,
		ProjectModule: projectModule,
	}

	// handler
	if err := generateModuleFile(data, "handler.go.tmpl",
		filepath.Join("handler", moduleName+".go")); err != nil {
		return err
	}

	// types (model + dto)
	if err := generateModuleFile(data, "model.go.tmpl",
		filepath.Join("types", "models", moduleName+".go")); err != nil {
		return err
	}
	if err := generateModuleFile(data, "dto.go.tmpl",
		filepath.Join("types", moduleName+".go")); err != nil {
		return err
	}

	// repository
	if err := generateModuleFile(data, "repo.go.tmpl",
		filepath.Join("repository", moduleName+".go")); err != nil {
		return err
	}

	// service
	if err := generateModuleFile(data, "service.go.tmpl",
		filepath.Join("service", moduleName+".go")); err != nil {
		return err
	}

	return nil
}

func generateModuleFile(data ModuleData, templateName, outputPath string) error {
	templateContent, err := moduleTemplateFS.ReadFile("templates/module-tmpl/" + templateName)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}
