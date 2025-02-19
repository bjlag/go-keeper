package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
)

const (
	configPathDefault = "./config/migrator.yaml"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", configPathDefault, "Configuration directory")
	flag.Parse()

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Environment: ", cfg.Env)
	fmt.Println("Migrations directory: ", cfg.SourcePath)
	fmt.Println("Migrations table name: ", cfg.MigrationsTable)
	fmt.Println("Database host: ", cfg.Database.Host)
	fmt.Println("Database port: ", cfg.Database.Port)
	fmt.Println("Database name: ", cfg.Database.Name)
	fmt.Println("Database user: ", cfg.Database.User)

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	db := pg.New(dsn).Connect()

	driver, err := pgx.WithInstance(db.DB, &pgx.Config{
		MigrationsTable: cfg.MigrationsTable,
	})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.SourcePath),
		cfg.Database.Name,
		driver,
	)
	if err != nil {
		log.Panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes")
			return
		}
		log.Panic(err)
	}

	fmt.Println("Migrations applied")
}
