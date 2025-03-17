// Отвечает за запуск сервера.
//
// Конфигурация указывается через флаг -c, описывается в YAML файле:
//   - пример ./config/server.yaml.dist
package main

import (
	"context"
	"flag"
	"fmt"
	logNative "log"
	"net"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/app/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
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

	dbCfg := cfg.Database
	db, err := pg.New(pg.GetDSN(dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.User, dbCfg.Password)).Connect()
	if err != nil {
		log.Error("Failed to get db connection", zap.Error(err))
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	jwt := auth.NewJWT(cfg.Auth.SecretKey, cfg.Auth.AccessTokenExp, cfg.Auth.RefreshTokenExp)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address.Host, cfg.Address.Port))
	if err != nil {
		log.Error("Failed to listen", zap.Error(err))
		panic(err)
	}

	if err = server.NewApp(db, jwt, listener, log).Run(ctx); err != nil {
		panic(err)
	}
}
