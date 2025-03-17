package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/jmoiron/sqlx"
)

type DBType string

// Константы содержат поддерживаемые базы данных.
const (
	TypePG     DBType = "pg"
	TypeSqlite        = "sqlite"
)

// Get возвращает настроенный экземпляр мигратора.
func Get(db *sqlx.DB, dbType DBType, dbName, sourcePath, migrationsTable string) (*migrate.Migrate, error) {
	const op = "migrator.Init"

	driver, err := getDBDriver(db, dbType, migrationsTable)
	if err != nil {
		return nil, fmt.Errorf("%s get driver: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", sourcePath),
		dbName,
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("%s create migrator instance: %w", op, err)
	}

	return m, nil
}

func getDBDriver(db *sqlx.DB, dbType DBType, migrationsTable string) (database.Driver, error) {
	switch dbType {
	case TypePG:
		return pgx.WithInstance(db.DB, &pgx.Config{
			MigrationsTable: migrationsTable,
		})
	case TypeSqlite:
		return sqlite3.WithInstance(db.DB, &sqlite3.Config{
			MigrationsTable: migrationsTable,
		})
	default:
		return nil, fmt.Errorf("uknown database type: %s", dbType) //nolint:err113
	}
}
