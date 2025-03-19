package backup

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

const prefixOp = "store.backup."

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Save(ctx context.Context, items []client.Backup) error {
	const op = prefixOp + "Save"

	rows := make([]row, 0, len(items))
	for _, item := range items {
		rows = append(rows, fromModel(item))
	}

	query := `
		INSERT INTO backup (guid, value) VALUES (:guid, :value)
		ON CONFLICT (guid) DO UPDATE
			SET value = excluded.value;
	`

	if _, err := s.db.NamedExecContext(ctx, query, rows); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Store) Get(ctx context.Context) ([]client.Backup, error) {
	const op = prefixOp + "Get"

	query := `SELECT guid, value FROM backup`

	var rows []row
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toModels(rows), nil
}

func (s *Store) Erase(ctx context.Context) error {
	const op = prefixOp + "Erase"

	query := `DELETE FROM backup;`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
