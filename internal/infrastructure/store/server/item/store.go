package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/bjlag/go-keeper/internal/domain/data"
)

var ErrNotAffectedRows = errors.New("not affected")

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

	var rows []row
	if err := s.db.SelectContext(ctx, &rows, query, userGUID, limit, offset); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return convertToModels(rows), nil
}

func (s *Store) Update(ctx context.Context, guid uuid.UUID, userGUID uuid.UUID, data model.UpdatedItem) error {
	const op = "store.item.Update"

	query := `
		UPDATE items
		SET encrypted_data = :encrypted_data, updated_at = :updated_at
		WHERE guid = :guid AND user_guid = :user_guid
	`

	arg := updated{
		GUID:          guid,
		UserGUID:      userGUID,
		EncryptedData: data.EncryptedData,
		UpdatedAt:     data.UpdatedAt,
	}

	result, err := s.db.NamedExecContext(ctx, query, arg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return ErrNotAffectedRows
	}

	return nil
}
