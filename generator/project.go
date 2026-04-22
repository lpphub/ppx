package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
	templates, err := discoverProjectTemplates()
	if err != nil {
		return fmt.Errorf("failed to discover templates: %w", err)
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

func discoverProjectTemplates() (map[string]string, error) {
	templates := make(map[string]string)

	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		dir := filepath.Dir(path)
		if dir == "templates/modules" {
			return nil
		}

		outputPath := templateToOutputPath(path)
		templates[path] = outputPath

		return nil
	})

	return templates, err
}

func templateToOutputPath(templatePath string) string {
	relPath := strings.TrimPrefix(templatePath, "templates/")
	outputPath := strings.TrimSuffix(relPath, ".tmpl")

	baseName := filepath.Base(outputPath)
	dir := filepath.Dir(outputPath)

	dotfiles := map[string]bool{
		"gitignore":   true,
		"env.example": true,
	}

	if dotfiles[baseName] && dir == "." {
		return "." + baseName
	}

	return outputPath
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
