package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// ValidateProjectName 验证项目名称是否合法
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// 检查长度
	if len(name) > 50 {
		return fmt.Errorf("project name too long (max 50 characters)")
	}

	// 检查是否只包含合法字符
	match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9-_]*$`, name)
	if !match {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, hyphens and underscores")
	}

	return nil
}

// GetGitConfig 获取git配置值
func GetGitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--global", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}