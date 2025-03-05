package item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/bjlag/go-keeper/internal/domain/data"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAllByUser(ctx context.Context, userGUID uuid.UUID, limit, offset uint32) ([]model.Item, error) {
	const op = "store.item.GetAllByUser"

	query := `
		SELECT guid, user_guid, encrypted_data, created_at, updated_at
		FROM items
		WHERE user_guid = $1
		ORDER BY guid
		LIMIT $2
		OFFSET $3
	`

	var rows []Row
	if err := s.db.SelectContext(ctx, &rows, query, userGUID, limit, offset); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return convertToModels(rows), nil
}
