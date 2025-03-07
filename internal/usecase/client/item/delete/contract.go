package delete

import (
	"context"

	"github.com/google/uuid"
)

type store interface {
	DeleteItem(ctx context.Context, guid uuid.UUID) error
}
