package update

import (
	"context"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/data"
)

type store interface {
	Update(ctx context.Context, guid uuid.UUID, userGUID uuid.UUID, updatedData model.UpdatedItem) error
}
