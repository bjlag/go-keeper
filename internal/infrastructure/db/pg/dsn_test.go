package pg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
)

func TestGetDSN(t *testing.T) {
	got := pg.GetDSN("host", "1111", "database_name", "user", "pass")

	assert.Equal(t, "postgresql://user:pass@host:1111/database_name?sslmode=disable", got)
}
