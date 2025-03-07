package item

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Store) SaveItem(ctx context.Context, item model.Item) error {
	const op = prefixOp + "SaveItem"

	query := `
		UPDATE items
		SET title = :title,
			value = :value,
			notes = :notes,
			updated_at = datetime()
		WHERE guid = :guid 	
	`

	r, err := toRow(item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.NamedExecContext(ctx, query, r)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// todo пределать на модель model.Item
func (s *Store) SaveItems(ctx context.Context, items []model.RawItem) error {
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
			INSERT INTO items (guid, category_id, title, value, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (guid) DO UPDATE SET
    			title = excluded.title,
    			value = excluded.value,
    			notes = excluded.notes,
    			updated_at = excluded.updated_at;
		`

		_, err := tx.ExecContext(ctx, query, i.GUID, i.Category, i.Title, i.Value, i.Notes, i.CreatedAt, i.UpdatedAt)
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

func (s *Store) ItemsByCategory(ctx context.Context, category model.Category) ([]model.RawItem, error) {
	const op = prefixOp + "Passwords"

	query := `
		SELECT guid, category_id, title, value, notes, created_at, updated_at
		FROM items
		WHERE category_id = $1;
	`
	var rows []row
	err := s.db.SelectContext(ctx, &rows, query, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return toModels(rows), nil
}
