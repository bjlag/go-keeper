package get_by_guid

import (
	"context"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

type dataStore interface {
	UserItemByGUID(ctx context.Context, userGUID, itemGUID uuid.UUID) (*model.Item, error)
}
