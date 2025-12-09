package generator

import (
	"io/fs"
	"os"
	"path/filepath"

	"ppx/internal/template"
	"ppx/internal/variables"
)

type Generator struct {
	dryRun bool
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(vars *variables.TemplateVars, targetDir string) error {
	projectPath := filepath.Join(targetDir, vars.ProjectName)

	// 创建项目目录
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	// 遍历模板文件并生成
	return fs.WalkDir(template.TemplateFS, "go-tpl", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel("go-tpl", path)
		if err != nil {
			return err
		}

		// 读取模板内容
		content, err := template.TemplateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// 处理变量替换
		replacer := variables.NewReplacer(vars)
		processed, err := replacer.Replace(content)
		if err != nil {
			return err
		}

		// 写入目标文件
		targetPath := filepath.Join(projectPath, relPath)

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// 处理特殊文件名
		if filepath.Base(targetPath) == "go.mod.template" {
			targetPath = filepath.Join(filepath.Dir(targetPath), "go.mod")
		}

		return os.WriteFile(targetPath, processed, 0644)
	})
}