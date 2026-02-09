package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"ppx/generator"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go web project",
	Long: `Create a new Go web project with predefined directory structure including cmd, infra, logic, web, and config folders.

Examples:
  ppx new myapp
  ppx new myapp --module github.com/user/myapp

Generated Project Structure:
  myapp/
  ├── config/
  │   └── config.yml
  ├── infra/
  │   ├── config.go
  │   ├── dbs.go
  │   ├── init.go
  │   └── jwt/
  │       └── jwt.go
  ├── logic/
  │   ├── dto/
  │   ├── auth/
  │   ├── user/
  │   ├── init.go
  │   └── wire.go
  ├── web/
  │   ├── middleware/
  │   ├── rest/
  │   └── app.go
  ├── main.go
  ├── go.mod
  └── Dockerfile`,
	Args: cobra.ExactArgs(1),
	Run:  runNew,
}

func runNew(cmd *cobra.Command, args []string) {
	projectName := args[0]

	// 验证项目名称
	if err := validateProjectName(projectName); err != nil {
		color.Red("❌ 无效的项目名称: %v", err)
		color.Yellow("💡 项目名称只能包含字母、数字、连字符和下划线")
		os.Exit(1)
	}

	// 检查目录是否已存在
	if _, err := os.Stat(projectName); err == nil {
		color.Red("❌ 目录 '%s' 已存在", projectName)
		color.Yellow("💡 请选择不同的名称或删除现有目录")
		os.Exit(1)
	}

	moduleName, _ := cmd.Flags().GetString("module")
	if moduleName == "" {
		moduleName = projectName
	}

	if err := generator.CreateProject(projectName, moduleName); err != nil {
		handleCreateError(err, projectName)
		os.Exit(1)
	}
}

func handleCreateError(err error, projectName string) {
	color.Red("❌ 创建项目 '%s' 失败", projectName)

	switch {
	case strings.Contains(err.Error(), "permission denied"):
		color.Yellow("💡 请尝试使用不同的权限或选择其他目录")
	case strings.Contains(err.Error(), "template"):
		color.Yellow("💡 这可能是模板中的错误，请报告此问题。")
	case strings.Contains(err.Error(), "disk space"):
		color.Yellow("💡 请检查可用磁盘空间")
	default:
		color.Yellow("💡 错误详情: %v", err)
	}

	color.Cyan("📞 需要帮助？访问: https://github.com/lpphub/ppx/issues")
}

func validateProjectName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("项目名称不能为空")
	}
	if len(name) > 50 {
		return fmt.Errorf("项目名称过长（最多50个字符）")
	}
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]*$`).MatchString(name) {
		return fmt.Errorf("项目名称必须以字母开头，且只能包含字母、数字、连字符和下划线")
	}
	return nil
}

func init() {
	newCmd.Flags().String("module", "", "项目的模块名（例如：github.com/user/project）")
}
