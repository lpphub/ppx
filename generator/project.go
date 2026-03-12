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

//go:embed templates/**
var templateFS embed.FS

type TemplateData struct {
	ProjectName string
	ModuleName  string
}

func CreateProject(projectName, moduleName string) error {
	bar := progressbar.NewOptions64(
		100,
		progressbar.OptionSetDescription(color.CyanString("🚀 Creating project...")),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	color.Cyan("\n📁 Creating directory structure...")
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	if err := createDirectories(projectName); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	_ = bar.Add(20)

	templateData := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	color.Cyan("📝 Processing template files...")
	if err := processTemplates(projectName, templateData, bar); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}

	_ = bar.Finish()
	printSuccess(projectName)
	return nil
}

func createDirectories(projectName string) error {
	directories := []string{
		"config",
		"contract",
		"infra/jwt",
		"module/auth",
		"module/user",
		"module/post",
		"server/helper",
		"server/middleware",
		"shared/consts",
		"shared/errs",
		"shared/mod",
		"shared/strutils",
	}

	for _, dir := range directories {
		dirPath := filepath.Join(projectName, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func processTemplates(projectName string, data TemplateData, bar *progressbar.ProgressBar) error {
	templates := map[string]string{
		"templates/static/main.go.tmpl":                     "main.go",
		"templates/static/go.mod.tmpl":                      "go.mod",
		"templates/static/Makefile.tmpl":                    "Makefile",
		"templates/static/Dockerfile.tmpl":                  "Dockerfile",
		"templates/static/gitignore.tmpl":                   ".gitignore",
		"templates/static/env.example.tmpl":                 ".env.example",
		"templates/dynamic/config/config.yml.tmpl":          "config/config.yml",
		"templates/dynamic/contract/user.go.tmpl":           "contract/user.go",
		"templates/dynamic/infra/init.go.tmpl":              "infra/init.go",
		"templates/dynamic/infra/config.go.tmpl":            "infra/config.go",
		"templates/dynamic/infra/database.go.tmpl":          "infra/database.go",
		"templates/dynamic/infra/jwt/jwt.go.tmpl":           "infra/jwt/jwt.go",
		"templates/dynamic/server/app.go.tmpl":              "server/app.go",
		"templates/dynamic/server/helper/helper.go.tmpl":    "server/helper/helper.go",
		"templates/dynamic/server/middleware/auth.go.tmpl":  "server/middleware/auth.go",
		"templates/dynamic/server/middleware/cors.go.tmpl":  "server/middleware/cors.go",
		"templates/dynamic/shared/consts/constants.go.tmpl": "shared/consts/constants.go",
		"templates/dynamic/shared/errs/errors.go.tmpl":      "shared/errs/errors.go",
		"templates/dynamic/shared/mod/module.go.tmpl":       "shared/mod/module.go",
		"templates/dynamic/shared/strutils/string.go.tmpl":  "shared/strutils/string.go",
		"templates/dynamic/module/user/init.go.tmpl":        "module/user/init.go",
		"templates/dynamic/module/user/model.go.tmpl":       "module/user/model.go",
		"templates/dynamic/module/user/dto.go.tmpl":         "module/user/dto.go",
		"templates/dynamic/module/user/handler.go.tmpl":     "module/user/handler.go",
		"templates/dynamic/module/user/service.go.tmpl":     "module/user/service.go",
		"templates/dynamic/module/user/repository.go.tmpl":  "module/user/repository.go",
		"templates/dynamic/module/auth/init.go.tmpl":        "module/auth/init.go",
		"templates/dynamic/module/auth/dto.go.tmpl":         "module/auth/dto.go",
		"templates/dynamic/module/auth/handler.go.tmpl":     "module/auth/handler.go",
		"templates/dynamic/module/auth/service.go.tmpl":     "module/auth/service.go",
		"templates/dynamic/module/post/init.go.tmpl":        "module/post/init.go",
		"templates/dynamic/module/post/model.go.tmpl":       "module/post/model.go",
		"templates/dynamic/module/post/dto.go.tmpl":         "module/post/dto.go",
		"templates/dynamic/module/post/handler.go.tmpl":     "module/post/handler.go",
		"templates/dynamic/module/post/service.go.tmpl":     "module/post/service.go",
		"templates/dynamic/module/post/repository.go.tmpl":  "module/post/repository.go",
	}

	templateCount := len(templates)
	progressPerFile := 80.0 / float64(templateCount)

	for templatePath, outputPath := range templates {
		if err := processTemplate(
			templatePath,
			filepath.Join(projectName, outputPath),
			data,
		); err != nil {
			return fmt.Errorf("failed to process template %s: %w", templatePath, err)
		}
		_ = bar.Add(int(progressPerFile))
	}

	return nil
}

func processTemplate(templatePath, outputPath string, data interface{}) error {
	templateContent, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}

func printSuccess(projectName string) {
	color.Green("\n🎉 Project '%s' created successfully!", projectName)

	color.Cyan("\n📂 Generated directory structure:")
	fmt.Printf("   %s/\n", projectName)
	fmt.Printf("   ├── config/\n")
	fmt.Printf("   ├── contract/\n")
	fmt.Printf("   ├── infra/\n")
	fmt.Printf("   │   └── jwt/\n")
	fmt.Printf("   ├── module/\n")
	fmt.Printf("   │   ├── auth/      # Authentication module\n")
	fmt.Printf("   │   ├── user/      # User module\n")
	fmt.Printf("   │   └── post/      # Demo CRUD module\n")
	fmt.Printf("   ├── server/\n")
	fmt.Printf("   │   ├── helper/\n")
	fmt.Printf("   │   └── middleware/\n")
	fmt.Printf("   └── shared/\n")
	fmt.Printf("       ├── consts/\n")
	fmt.Printf("       ├── errs/\n")
	fmt.Printf("       ├── mod/\n")
	fmt.Printf("       └── strutils/\n")

	color.Cyan("\n📋 Next steps:")
	fmt.Printf("   1. cd %s\n", projectName)
	fmt.Printf("   2. Update config/config.yml with your database credentials\n")
	fmt.Printf("   3. cp .env.example .env && edit .env for local development\n")
	fmt.Printf("   4. go mod tidy\n")
	fmt.Printf("   5. go run .\n")

	color.Yellow("\n⚠ Don't forget:")
	fmt.Printf("   - Update config/config.yml (Database, Redis, JWT settings)\n")
	fmt.Printf("   - Default server port: 8080\n")

	color.Cyan("\n📚 Documentation: https://github.com/lpphub/ppx\n")
}
