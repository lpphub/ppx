package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	ExpireTime       int64  `mapstructure:"expire_time"`
	RefreshExpireTime int64  `mapstructure:"refresh_expire_time"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

var Cfg *Config

func Init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		panic("Failed to unmarshal config: " + err.Error())
	}
}