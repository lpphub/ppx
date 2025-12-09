package infra

import (
	"{{.ModulePath}}/infra/config"
	"{{.ModulePath}}/infra/dbs"
	"{{.ModulePath}}/infra/jwt"
	"{{.ModulePath}}/infra/logger"
	"{{.ModulePath}}/infra/monitor"
)

func Init() {
	config.Init()
	logger.Init()
	dbs.Init()
	monitor.Init()
	// jwt不需要初始化，只是工具函数
}