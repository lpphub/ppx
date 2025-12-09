package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"ppx/internal/generator"
	"ppx/internal/variables"
	"ppx/internal/utils"
)

var (
	modulePath    string
	authorName    string
	authorEmail   string
	description   string
	databaseType  string
	noRedis       bool
	noMetrics     bool
	noPprof       bool
	targetDir     string
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new Go web project",
	Long: `Create a new Go web project with clean architecture template.
The project will include user management, RBAC, JWT authentication, and other production-ready features.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// 基本配置
	createCmd.Flags().StringVar(&modulePath, "module-path", "", "Go module path (default: github.com/user/project)")
	createCmd.Flags().StringVar(&authorName, "author", "", "Author name (default: git config user.name)")
	createCmd.Flags().StringVar(&authorEmail, "email", "", "Author email (default: git config user.email)")
	createCmd.Flags().StringVar(&description, "description", "", "Project description")

	// 技术选项
	createCmd.Flags().StringVar(&databaseType, "database", "mysql", "Database type (mysql, postgres, sqlite)")
	createCmd.Flags().BoolVar(&noRedis, "no-redis", false, "Disable Redis cache")
	createCmd.Flags().BoolVar(&noMetrics, "no-metrics", false, "Disable Prometheus metrics")
	createCmd.Flags().BoolVar(&noPprof, "no-pprof", false, "Disable pprof profiling")

	// 输出选项
	createCmd.Flags().StringVar(&targetDir, "target", ".", "Target directory")
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// 创建模板变量
	vars := variables.NewDefaults(projectName)

	// 覆盖默认值
	if modulePath != "" {
		vars.ModulePath = modulePath
	}
	if authorName != "" {
		vars.AuthorName = authorName
	} else {
		// 从git config获取
		if name, _ := utils.GetGitConfig("user.name"); name != "" {
			vars.AuthorName = name
		}
	}
	if authorEmail != "" {
		vars.AuthorEmail = authorEmail
	} else {
		// 从git config获取
		if email, _ := utils.GetGitConfig("user.email"); email != "" {
			vars.AuthorEmail = email
		}
	}
	if description != "" {
		vars.Description = description
	}

	// 设置功能开关
	vars.DatabaseType = databaseType
	vars.RedisEnabled = !noRedis
	vars.MetricsEnabled = !noMetrics
	vars.PprofEnabled = !noPprof

	// 验证输入
	if err := utils.ValidateProjectName(projectName); err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// 检查目标目录
	projectPath := filepath.Join(targetDir, projectName)
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists", projectPath)
	}

	// 生成项目
	fmt.Printf("Creating project %s...\n", projectName)
	g := generator.New()
	if err := g.Generate(vars, targetDir); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Printf("\nProject %s created successfully!\n", projectName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run .\n")

	return nil
}