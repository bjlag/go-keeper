package logger

import (
	"os"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *zap.Logger
)

func Get(env string) *zap.Logger {
	once.Do(func() {
		var config zap.Config

		if env == "prod" {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
		}

		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.InitialFields = map[string]interface{}{
			"env":        env,
			"pid":        os.Getpid(),
			"go_version": runtime.Version(),
		}

		logger = zap.Must(config.Build())
	})

	return logger
}
