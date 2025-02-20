package main

import (
	"errors"
	"flag"
	"fmt"
	logNative "log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
)

const (
	configPathDefault = "./config/migrator.yaml"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logNative.Fatalf("panic occurred: %v", r)
		}
	}()

	var configPath string

	flag.StringVar(&configPath, "c", configPathDefault, "Configuration directory")
	flag.Parse()

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Env)
	defer func() {
		_ = log.Sync()
	}()

	log.Debug("Config loaded", zap.Any("config", cfg))

	dbConf := cfg.Database
	db, err := pg.New(pg.GetDSN(dbConf.Host, dbConf.Port, dbConf.Name, dbConf.User, dbConf.Password)).Connect()
	if err != nil {
		panic(err)
	}

	m, err := initMigrator(db, cfg)
	if err != nil {
		log.Error("Failed to initialize migrator", zap.Error(err))
		panic(err)
	}

	log.Info("Applying migrations")

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No changes")
			return
		}
		log.Error("Migration failed", zap.Error(err))
		return
	}

	log.Info("Migrations applied")
}

func initMigrator(db *sqlx.DB, cfg Config) (*migrate.Migrate, error) {
	const op = "initMigrator"

	driver, err := pgx.WithInstance(db.DB, &pgx.Config{
		MigrationsTable: cfg.MigrationsTable,
	})
	if err != nil {
		return nil, fmt.Errorf("%s get driver: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.SourcePath),
		cfg.Database.Name,
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("%s create migrator instance: %w", op, err)
	}

	return m, nil
}
