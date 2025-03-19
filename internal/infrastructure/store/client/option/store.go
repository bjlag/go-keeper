// Package option отвечает за работу с опциями в базе данных на стороне клиента.
package option

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

const prefixOp = "store.option."

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// SaveOption сохраняет переданную модель опции в базе данных.
// Если элемента в базе нет, то создает.
func (s *Store) SaveOption(ctx context.Context, option model.Option) error {
	const op = prefixOp + "SaveOption"

	query := `
		INSERT INTO options (slug, value)
		VALUES (:slug, :value)
		ON CONFLICT (slug) DO UPDATE 
			SET value = excluded.value 	
	`

	r := toRow(option)
	_, err := s.db.NamedExecContext(ctx, query, r)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// OptionBySlug получает опцию по ее слагу.
func (s *Store) OptionBySlug(ctx context.Context, slug string) (*model.Option, error) {
	const op = prefixOp + "OptionBySlug"

	query := `
		SELECT slug, value
		FROM options
		WHERE slug = $1;
	`
	var r row
	err := s.db.GetContext(ctx, &r, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	m := r.toModel()

	return &m, nil
}
