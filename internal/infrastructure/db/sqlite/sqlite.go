package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	dsn string
}

func New(dsn string) *SQLite {
	return &SQLite{
		dsn: dsn,
	}
}

func (l SQLite) Connect() (*sqlx.DB, error) {
	const op = "sqlite.Connect"

	db, err := sqlx.Connect("sqlite3", l.dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
