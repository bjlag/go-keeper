package item

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

const prefixOp = "store.item."

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) SaveItems(ctx context.Context, items []model.Item) error {
	const op = prefixOp + "SaveItems"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, i := range items {
		query := `
			INSERT INTO items (guid, data, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (guid) DO UPDATE SET
    			data = excluded.data,
    			updated_at = excluded.updated_at;
		`

		_, err := tx.ExecContext(ctx, query, i.GUID, i.Data, i.CreatedAt, i.UpdatedAt)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
