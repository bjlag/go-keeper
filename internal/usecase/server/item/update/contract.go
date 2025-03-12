package update

import (
	"context"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

type store interface {
	ItemByGUID(ctx context.Context, guid uuid.UUID) (*model.Item, error)
	Update(ctx context.Context, guid uuid.UUID, userGUID uuid.UUID, updatedData model.UpdatedItem) error
}
