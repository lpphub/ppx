package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

//go:embed templates/*
var templateFS embed.FS

// TemplateData holds the data for template rendering
type TemplateData struct {
	ProjectName string
	ModuleName  string
}

// CreateProject creates a new Go web project with the given name and module
func CreateProject(projectName, moduleName string) error {
	// 创建进度条
	bar := progressbar.NewOptions64(
		100,
		progressbar.OptionSetDescription("Creating project..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			color.Green("✓ Project created successfully!")
		}),
	)

	// 步骤1: 创建项目根目录
	color.Cyan("📁 Creating directory structure...")
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}
	_ = bar.Add(10)

	// 步骤2: 创建目录结构
	if err := createDirectories(projectName); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	_ = bar.Add(20)

	// 准备模板数据
	if moduleName == "" {
		moduleName = projectName
	}
	templateData := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	// 步骤3: 处理模板
	color.Cyan("📝 Processing templates...")
	if err := processTemplates(projectName, templateData, bar); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}
	_ = bar.Add(70)

	_ = bar.Finish()

	printNextSteps(projectName)
	return nil
}

func printNextSteps(projectName string) {
	color.Green("\n🎉 Project '%s' created successfully!", projectName)
	color.Cyan("\n📋 Next steps:")
	fmt.Printf("   1. cd %s\n", projectName)
	fmt.Printf("   2. Update config/config.yml with your database settings\n")
	fmt.Printf("   3. go mod download\n")
	fmt.Printf("   4. wire ./logic  # Generate dependency injection code\n")
	fmt.Printf("   5. go run .\n")

	color.Yellow("\n⚠ Don't forget to:")
	fmt.Printf("   - Update database credentials in config/config.yml\n")
	fmt.Printf("   - Change JWT secret in config/config.yml\n")
	fmt.Printf("   - Start Redis server if needed\n")
	fmt.Printf("   - Install wire if not available: go install github.com/google/wire/cmd/wire@latest\n")
}

// createDirectories creates the required directory structure for the project
func createDirectories(projectName string) error {
	// Define directories to create
	directories := []string{
		"infra/jwt",
		"logic/auth",
		"logic/user",
		"logic/shared",
		"web/middleware",
		"web/rest",
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

	return nil
}

func processTemplates(projectName string, data TemplateData, bar *progressbar.ProgressBar) error {
	// All templates mapping
	templates := map[string]string{
		"templates/go.mod.tmpl":                     "go.mod",
		"templates/Dockerfile.tmpl":                 "Dockerfile",
		"templates/main.go.tmpl":                    "main.go",
		"templates/config/config.yml.tmpl":          "config/config.yml",
		"templates/infra/init.go.tmpl":              "infra/init.go",
		"templates/infra/config.go.tmpl":            "infra/config.go",
		"templates/infra/db.go.tmpl":                "infra/db.go",
		"templates/infra/jwt/jwt.go.tmpl":           "infra/jwt/jwt.go",
		"templates/logic/auth/service.go.tmpl":      "logic/auth/service.go",
		"templates/logic/user/service.go.tmpl":      "logic/user/service.go",
		"templates/logic/user/model.go.tmpl":        "logic/user/model.go",
		"templates/logic/shared/consts.go.tmpl":     "logic/shared/consts.go",
		"templates/logic/shared/errors.go.tmpl":     "logic/shared/errors.go",
		"templates/logic/shared/pagination.go.tmpl": "logic/shared/pagination.go",
		"templates/logic/init.go.tmpl":              "logic/init.go",
		"templates/logic/wire.go.tmpl":              "logic/wire.go",
		"templates/web/app.go.tmpl":                 "web/app.go",
		"templates/web/middleware/auth.go.tmpl":     "web/middleware/auth.go",
		"templates/web/rest/handler.go.tmpl":        "web/rest/handler.go",
		"templates/web/rest/user/handler.go.tmpl":   "web/rest/user/handler.go",
		"templates/web/rest/user/route.go.tmpl":     "web/rest/user/route.go",
		"templates/web/types/req.go.tmpl":           "web/types/req.go",
		"templates/web/types/resp.go.tmpl":          "web/types/resp.go",
	}

	// Process all templates
	templateCount := len(templates)
	for i, templatePath := range []string{
		"templates/go.mod.tmpl",
		"templates/Dockerfile.tmpl",
		"templates/main.go.tmpl",
		"templates/config/config.yml.tmpl",
		"templates/infra/init.go.tmpl",
		"templates/infra/config.go.tmpl",
		"templates/infra/db.go.tmpl",
		"templates/infra/jwt/jwt.go.tmpl",
		"templates/logic/auth/service.go.tmpl",
		"templates/logic/user/service.go.tmpl",
		"templates/logic/user/model.go.tmpl",
		"templates/logic/shared/consts.go.tmpl",
		"templates/logic/shared/errors.go.tmpl",
		"templates/logic/shared/pagination.go.tmpl",
		"templates/logic/init.go.tmpl",
		"templates/logic/wire.go.tmpl",
		"templates/web/app.go.tmpl",
		"templates/web/middleware/auth.go.tmpl",
		"templates/web/rest/handler.go.tmpl",
		"templates/web/rest/user/handler.go.tmpl",
		"templates/web/rest/user/route.go.tmpl",
		"templates/web/types/req.go.tmpl",
		"templates/web/types/resp.go.tmpl",
	} {
		outputPath := templates[templatePath]
		if err := processTemplate(
			templatePath,
			filepath.Join(projectName, outputPath),
			data,
		); err != nil {
			return fmt.Errorf("failed to process template %s: %w", templatePath, err)
		}

		// 更新进度条
		progress := int(float64(i+1) / float64(templateCount) * 70)
		_ = bar.Set(progress)
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
