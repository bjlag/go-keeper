package remove

import (
	"context"

	"github.com/google/uuid"
)

type store interface {
	Delete(ctx context.Context, guid uuid.UUID, userGUID uuid.UUID) error
}
