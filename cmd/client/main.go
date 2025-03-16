// Отвечает за запуск клиента.
//
// Конфигурация указывается через флаг -c, описывается в YAML файле:
//   - пример ./config/client.yaml.dist
//
// Флаг -version выведет текущую версию и дату сборки.
package main

import (
	"context"
	"flag"
	"fmt"
	logNative "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/app/client"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
)

const (
	// Конфигурация по умолчанию, если не передать флаг конфигурации при старте приложения.
	configPathDefault = "./config/client.yaml"
)

var (
	viewVersion  bool
	buildVersion string
	buildDate    string
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logNative.Fatalf("panic occurred: %v", r)
		}
	}()

	var configPath string

	flag.StringVar(&configPath, "c", configPathDefault, "Path to config file")
	flag.BoolVar(&viewVersion, "version", false, "View build version and data")
	flag.Parse()

	if viewVersion {
		fmt.Printf("Version: %s\nBuild: %s\n", buildVersion, buildDate)
		os.Exit(0)
	}

	var cfg client.Config
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

	if err := client.NewApp(cfg, log).Run(ctx); err != nil {
		panic(err)
	}
}
