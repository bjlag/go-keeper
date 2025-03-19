package get_all

import (
	"context"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

type dataStore interface {
	GetAllByUser(ctx context.Context, userGUID uuid.UUID, limit, offset uint32) ([]model.Item, error)
}
