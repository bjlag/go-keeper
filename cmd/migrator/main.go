// Отвечает за применение миграций к базе данных.
// Поддерживает PostgreSQL и SQLite.
//
// Конфигурация указывается через флаг -c, описывается в YAML файле:
//   - пример для клиента ./config/migrator_client.yaml.dist
//   - пример для сервера ./config/migrator_server.yaml.dist
package main

import (
	"errors"
	"flag"
	"io/fs"
	logNative "log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/sqlite"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/infrastructure/migrator"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logNative.Fatalf("panic occurred: %v", r)
		}
	}()

	var configPath string

	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.Parse()

	var cfg migrator.Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Env)
	defer func() {
		_ = log.Sync()
	}()

	log.Debug("Config loaded", zap.Any("config", cfg))
	log = log.With(zap.Any("db_type", cfg.Database.Type))

	var (
		err error
		db  *sqlx.DB
	)

	switch cfg.Database.Type {
	case migrator.TypePG:
		dbConf := cfg.Database
		db, err = pg.New(pg.GetDSN(dbConf.Host, dbConf.Port, dbConf.Name, dbConf.User, dbConf.Password)).Connect()
	case migrator.TypeSqlite:
		db, err = sqlite.New("./client.db").Connect()
	}
	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		log.Error("Failed to open db", zap.Error(err))
		panic(err)
	}

	m, err := migrator.Get(db, cfg.Database.Type, cfg.Database.Name, cfg.SourcePath, cfg.MigrationsTable)
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
		var e *fs.PathError
		if errors.As(err, &e) {
			log.Info("No migration files")
			return
		}
		log.Error("Migration failed", zap.Error(err))
		return
	}

	log.Info("Migrations applied")
}
