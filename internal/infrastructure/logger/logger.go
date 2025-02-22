package logger

import (
	"context"
	"os"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKeyLogger int

const IDLoggerKey ctxKeyLogger = 0

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

func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(IDLoggerKey).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(IDLoggerKey).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, IDLoggerKey, l)
}
