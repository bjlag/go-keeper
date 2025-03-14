// Package item отвечает за работу с элементами в базе данных на стороне клиента.
package item

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

// SaveItem сохраняет переданный элемент в базе.
func (s *Store) SaveItem(ctx context.Context, item model.Item) error {
	const op = prefixOp + "SaveItem"

	query := `
		UPDATE items
		SET title = :title,
			value = :value,
			notes = :notes,
			updated_at = :updated_at
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

// CreateItem создает переданную модель элемента в базе.
func (s *Store) CreateItem(ctx context.Context, item model.Item) error {
	const op = prefixOp + "CreateItem"

	query := `
		INSERT INTO items(guid, category_id, title, value, notes, updated_at, created_at)
		VALUES (:guid, :category_id, :title, :value, :notes, :updated_at, :created_at)
	`

	var value *[]byte

	if item.Value != nil {
		v, err := json.Marshal(item.Value)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		value = &v
	}

	args := row{
		GUID:      item.GUID,
		Category:  item.Category,
		Title:     item.Title,
		Value:     value,
		Notes:     item.Notes,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	_, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteItem удаляет элемент с переданным GUID из базы.
func (s *Store) DeleteItem(ctx context.Context, guid uuid.UUID) error {
	const op = prefixOp + "DeleteItem"

	query := `
		DELETE FROM items WHERE guid = $1
	`

	_, err := s.db.ExecContext(ctx, query, guid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// SaveItems сохраняет несколько элементов.
func (s *Store) SaveItems(ctx context.Context, items []model.RawItem) error {
	// todo пределать на модель model.Item
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

// ItemsByCategory получает элементы указанной категории.
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
