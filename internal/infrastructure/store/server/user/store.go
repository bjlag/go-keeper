// Package user отвечает за работу с пользователями на стороне сервера.
package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "github.com/bjlag/go-keeper/internal/domain/server/user"
)

// ErrNotFound ошибка в случае, если пользователь не найден.
var ErrNotFound = errors.New("user not found")

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// GetByGUID получить пользователя по его GUID.
func (s Store) GetByGUID(ctx context.Context, guid uuid.UUID) (*model.User, error) {
	const op = "store.user.GetByGUID"

	query := `
		SELECT guid, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE guid = $1
	`

	var user row
	if err := s.db.GetContext(ctx, &user, query, guid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user.convertToModel(), nil
}

// GetByEmail получить пользователя по его емейлу.
func (s Store) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "store.user.GetByEmail"

	query := `
		SELECT guid, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	var user row
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user.convertToModel(), nil
}

// Add добавить пользователя.
func (s Store) Add(ctx context.Context, user model.User) error {
	const op = "store.user.Add"

	query := `INSERT INTO users (guid, email, password_hash) VALUES (:guid, :email, :password_hash)`
	_, err := s.db.NamedExecContext(ctx, query, convertToRow(user))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
