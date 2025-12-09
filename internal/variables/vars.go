package variables

import "time"

type TemplateVars struct {
	// 基础信息
	ProjectName   string
	ModulePath    string
	Description   string

	// 作者信息
	AuthorName    string
	AuthorEmail   string

	// 功能开关
	DatabaseType  string
	RedisEnabled  bool
	MetricsEnabled bool
	PprofEnabled  bool

	// 时间信息
	GeneratedAt   string
}

func NewDefaults(projectName string) *TemplateVars {
	return &TemplateVars{
		ProjectName:    projectName,
		ModulePath:     "github.com/user/" + projectName,
		Description:    "A Go web application built with clean architecture",
		DatabaseType:   "mysql",
		RedisEnabled:   true,
		MetricsEnabled: true,
		PprofEnabled:   true,
		GeneratedAt:    time.Now().Format("2006-01-02"),
	}
}