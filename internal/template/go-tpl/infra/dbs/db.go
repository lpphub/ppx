package dbs

import (
	"fmt"
	"time"

	"{{.ModulePath}}/infra/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		config.Cfg.Database.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database instance: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}