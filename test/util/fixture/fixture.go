package fixture

import (
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
)

// Load загружает в БД по переданному подключению db фикстуры из директории dir.
func Load(db *sqlx.DB, dir string) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dir),
	)
	if err != nil {
		return err
	}

	if err := fixtures.Load(); err != nil {
		return err
	}

	return nil
}
