package infra

import (
	"fmt"

	"github.com/lpphub/goweb/pkg/logging"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Cfg *Config
	DB  *gorm.DB
	RDB *redis.Client
)

func Init() {
	var err error
	// 1.加载配置
	Cfg, err = LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2.配置日志
	logging.Init()

	// 3.初始化数据库和Redis
	DB, err = NewMysqlDB(Cfg.Database)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	RDB, err = NewRedis(Cfg.Redis)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %v", err))
	}

}

func ProvideDB() *gorm.DB {
	return DB
}

func ProvideRDB() *redis.Client {
	return RDB
}
