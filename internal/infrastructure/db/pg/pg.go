package pg

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func (p *PG) Connect() *sqlx.DB {
	db, err := sqlx.Connect("pgx", p.dsn)
	if err != nil {
		log.Panic(err)
	}

	db.SetMaxOpenConns(maxOpenConnects)
	db.SetMaxIdleConns(maxIdleConnects)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	return db
}
