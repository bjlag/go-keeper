package pg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
)

func TestPG_Connect(t *testing.T) {
	got, err := pg.New("bad_dsn").Connect()
	assert.Error(t, err) //nolint:testifylint
	assert.Nil(t, got)
}
