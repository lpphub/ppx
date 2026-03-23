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
		"infra/jwt",
		"modules/core",
		"modules/auth",
		"modules/user",
		"modules/post",
		"server/helper",
		"server/middleware",
		"shared/consts",
		"shared/contracts",
		"shared/errs",
		"shared/pagination",
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
		"templates/main.go.tmpl":                     "main.go",
		"templates/go.mod.tmpl":                      "go.mod",
		"templates/Makefile.tmpl":                    "Makefile",
		"templates/Dockerfile.tmpl":                  "Dockerfile",
		"templates/gitignore.tmpl":                   ".gitignore",
		"templates/env.example.tmpl":                 ".env.example",
		"templates/config/config.yml.tmpl":           "config/config.yml",
		"templates/infra/init.go.tmpl":               "infra/init.go",
		"templates/infra/config.go.tmpl":             "infra/config.go",
		"templates/infra/dbs.go.tmpl":                "infra/dbs.go",
		"templates/infra/jwt/jwt.go.tmpl":            "infra/jwt/jwt.go",
		"templates/server/app.go.tmpl":               "server/app.go",
		"templates/server/helper/response.go.tmpl":   "server/helper/response.go",
		"templates/server/middleware/auth.go.tmpl":   "server/middleware/auth.go",
		"templates/server/middleware/cors.go.tmpl":   "server/middleware/cors.go",
		"templates/shared/consts/constants.go.tmpl":  "shared/consts/constants.go",
		"templates/shared/contracts/user.go.tmpl":    "shared/contracts/user.go",
		"templates/shared/errs/errors.go.tmpl":       "shared/errs/errors.go",
		"templates/shared/pagination/cursor.go.tmpl": "shared/pagination/cursor.go",
		"templates/shared/pagination/offset.go.tmpl": "shared/pagination/offset.go",
		"templates/shared/strutils/string.go.tmpl":   "shared/strutils/string.go",
		"templates/modules/core/module.go.tmpl":      "modules/core/module.go",
		"templates/modules/user/module.go.tmpl":      "modules/user/module.go",
		"templates/modules/user/model.go.tmpl":       "modules/user/model.go",
		"templates/modules/user/dto.go.tmpl":         "modules/user/dto.go",
		"templates/modules/user/handler.go.tmpl":     "modules/user/handler.go",
		"templates/modules/user/service.go.tmpl":     "modules/user/service.go",
		"templates/modules/user/repo.go.tmpl":        "modules/user/repo.go",
		"templates/modules/auth/module.go.tmpl":      "modules/auth/module.go",
		"templates/modules/auth/dto.go.tmpl":         "modules/auth/dto.go",
		"templates/modules/auth/handler.go.tmpl":     "modules/auth/handler.go",
		"templates/modules/auth/service.go.tmpl":     "modules/auth/service.go",
		"templates/modules/post/module.go.tmpl":      "modules/post/module.go",
		"templates/modules/post/model.go.tmpl":       "modules/post/model.go",
		"templates/modules/post/dto.go.tmpl":         "modules/post/dto.go",
		"templates/modules/post/handler.go.tmpl":     "modules/post/handler.go",
		"templates/modules/post/service.go.tmpl":     "modules/post/service.go",
		"templates/modules/post/repo.go.tmpl":        "modules/post/repo.go",
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
	fmt.Printf("   │   └── config.yml\n")
	fmt.Printf("   ├── modules/\n")
	fmt.Printf("   │   ├── core/        # Module interface\n")
	fmt.Printf("   │   ├── auth/        # Authentication module\n")
	fmt.Printf("   │   ├── user/        # User module\n")
	fmt.Printf("   │   └── post/        # Demo CRUD module\n")
	fmt.Printf("   ├── infra/\n")
	fmt.Printf("   │   ├── dbs.go\n")
	fmt.Printf("   │   └── jwt/\n")
	fmt.Printf("   ├── server/\n")
	fmt.Printf("   │   ├── helper/\n")
	fmt.Printf("   │   └── middleware/\n")
	fmt.Printf("   ├── shared/\n")
	fmt.Printf("   │   ├── consts/\n")
	fmt.Printf("   │   ├── contracts/   # Module contracts\n")
	fmt.Printf("   │   ├── errs/\n")
	fmt.Printf("   │   ├── pagination/\n")
	fmt.Printf("   │   └── strutils/\n")
	fmt.Printf("   ├── main.go\n")
	fmt.Printf("   ├── go.mod\n")
	fmt.Printf("   ├── Makefile\n")
	fmt.Printf("   └── Dockerfile\n")

	color.Cyan("\n📋 Next steps:")
	fmt.Printf("   1. cd %s\n", projectName)
	fmt.Printf("   2. Update config/config.yml with your database credentials\n")
	fmt.Printf("   3. mv .env.example .env && edit .env for local development\n")
	fmt.Printf("   4. go mod tidy\n")
	fmt.Printf("   5. go run .\n")

	color.Yellow("\n⚠ Don't forget:")
	fmt.Printf("   - Update config/config.yml (Database, Redis, JWT settings)\n")
	fmt.Printf("   - Default server port: 8080\n")

	color.Cyan("\n📚 Documentation: https://github.com/lpphub/ppx\n")
}
