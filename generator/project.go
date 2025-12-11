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
		progressbar.OptionSetDescription("正在创建项目..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			color.Green("✓ 项目创建成功！")
		}),
	)

	// 步骤1: 创建项目根目录
	color.Cyan("📁 创建目录结构...")
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
	color.Cyan("📝 处理模板文件...")
	if err := processTemplates(projectName, templateData, bar); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}
	_ = bar.Add(70)

	_ = bar.Finish()

	printNextSteps(projectName)
	return nil
}

func printNextSteps(projectName string) {
	color.Green("\n🎉 项目 '%s' 创建成功！", projectName)
	color.Cyan("\n📋 接下来的步骤:")
	fmt.Printf("   1. cd %s\n", projectName)
	fmt.Printf("   2. 更新 config/config.yml 中的配置\n")
	fmt.Printf("   3. go mod tidy\n")
	fmt.Printf("   4. wire ./logic  # 生成依赖注入代码\n")
	fmt.Printf("   5. go run .\n")

	color.Yellow("\n⚠ 不要忘记:")
	fmt.Printf("   - 更新 config/config.yml 中的配置（DB、Redis、JWT）\n")
	fmt.Printf("   - 如果没有安装 wire: go install github.com/google/wire/cmd/wire@latest\n")
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
	i := 0
	for templatePath, outputPath := range templates {
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
		i++
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
