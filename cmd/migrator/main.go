package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
)

const (
	sourcePathDefault = "./migrations"
)

func main() {
	var sourcePath string

	flag.StringVar(&sourcePath, "dir", sourcePathDefault, "Migrations directory")
	flag.Parse()

	fmt.Println("Migrations directory: ", sourcePath)

	dsn := "postgresql://postgres:secret@localhost:5444/master?sslmode=disable"
	db := pg.New(dsn).Connect()

	driver, err := pgx.WithInstance(db.DB, &pgx.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", sourcePath),
		"master",
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
