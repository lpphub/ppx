package logger

import (
	"os"
	"{{.ModulePath}}/infra/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log zerolog.Logger

func Init() {
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

	Log = zerolog.New(multi).With().Timestamp().Logger()

	if config.Cfg.Server.Mode == "debug" {
		Log = Log.Level(zerolog.DebugLevel)
	} else {
		Log = Log.Level(zerolog.InfoLevel)
	}
}