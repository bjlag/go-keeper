// Package logger отвечает за создание логгера и работу с ним.
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

// IDLoggerKey ID ключа для работы с логгером через контекст.
const IDLoggerKey ctxKeyLogger = 0

var (
	once   sync.Once
	logger *zap.Logger
)

// Get получить логгер. Используется паттерн Singletone.
// При первом получении создается экземпляр логгера и кладется в глобальную переменную пакета.
func Get(env string) *zap.Logger {
	once.Do(func() {
		if env == "test" {
			logger = zap.NewNop()
			return
		}

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

// FromCtx получает логгер из контекста.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(IDLoggerKey).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx кладет логгер в контекст.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(IDLoggerKey).(*zap.Logger); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, IDLoggerKey, l)
}
