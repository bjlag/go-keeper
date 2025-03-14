// Отвечает за запуск сервера.
//
// Конфигурация указывается через флаг -c, описывается в YAML файле:
//   - пример ./config/server.yaml.dist
package main

import (
	"context"
	"flag"
	logNative "log"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/app/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
)

const (
	// Конфигурация по умолчанию, если не передать флаг конфигурации при старте приложения.
	configPathDefault = "./config/server.yaml"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logNative.Fatalf("panic occurred: %v", r)
		}
	}()

	var configPath string

	flag.StringVar(&configPath, "c", configPathDefault, "Path to config file")
	flag.Parse()

	var cfg server.Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Env)
	defer func() {
		_ = log.Sync()
	}()

	log.Debug("Config loaded", zap.Any("config", cfg))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	err := server.NewApp(cfg, log).Run(ctx)
	if err != nil {
		panic(err)
	}
}
