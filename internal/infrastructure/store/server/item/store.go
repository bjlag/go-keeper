// Package item отвечает за работу с элементами в базе данных на стороне сервера.
package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

var (
	// ErrNotAffectedRows ошибка в случае, если при обновлении не было задето ни одной записи.
	ErrNotAffectedRows = errors.New("not affected")
	// ErrNotFound ошибка в случае, если не было найдено ни одной записи.
	ErrNotFound = errors.New("not found")
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// GetAllByUser получить все элементы пользователя.
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

// ItemByGUID получить элемент по его GUID.
func (s *Store) ItemByGUID(ctx context.Context, guid uuid.UUID) (*model.Item, error) {
	const op = "store.item.ItemByGUID"

	query := `
		SELECT guid, user_guid, encrypted_data, created_at, updated_at
		FROM items
		WHERE guid = $1
	`

	var r row
	if err := s.db.GetContext(ctx, &r, query, guid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := r.convertToModel()

	return &result, nil
}

// UserItemByGUID получить элемент пользователя по его GUID.
func (s *Store) UserItemByGUID(ctx context.Context, userGUID, itemGUID uuid.UUID) (*model.Item, error) {
	const op = "store.item.UserItemByGUID"

	query := `
		SELECT guid, user_guid, encrypted_data, created_at, updated_at
		FROM items
		WHERE guid = $1 AND user_guid = $2
	`

	var r row
	if err := s.db.GetContext(ctx, &r, query, itemGUID, userGUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := r.convertToModel()

	return &result, nil
}

// Create создать элемент по переданной модели.
func (s *Store) Create(ctx context.Context, item model.Item) error {
	const op = "store.item.Create"

	query := `
		INSERT INTO items(guid, user_guid, encrypted_data, created_at, updated_at)
		VALUES(:guid, :user_guid, :encrypted_data, :created_at, :updated_at)
	`

	arg := row{
		GUID:          item.GUID,
		UserGUID:      item.UserGUID,
		EncryptedData: item.EncryptedData,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}

	_, err := s.db.NamedExecContext(ctx, query, arg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Update обновить элемент пользователя.
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

// Delete удалить элемент пользователя.
func (s *Store) Delete(ctx context.Context, guid uuid.UUID, userGUID uuid.UUID) error {
	const op = "store.item.Delete"

	query := `
		DELETE FROM items
		WHERE guid = :guid AND user_guid = :user_guid
	`

	arg := deleted{
		GUID:     guid,
		UserGUID: userGUID,
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
