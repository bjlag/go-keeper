package pg

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

const (
	maxOpenConnects = 5
	maxIdleConnects = 5
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 5 * time.Minute
)

type PG struct {
	dsn string
}

func New(dsn string) *PG {
	return &PG{
		dsn: dsn,
	}
}

func (p *PG) Connect() (*sqlx.DB, error) {
	const op = "pg.Connect"

	db, err := sqlx.Connect("pgx", p.dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(maxOpenConnects)
	db.SetMaxIdleConns(maxIdleConnects)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	return db, nil
}
