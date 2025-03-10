package delete

import (
	"context"

	"github.com/google/uuid"

	dto "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type rpc interface {
	DeleteItem(ctx context.Context, in *dto.DeleteItemIn) error
}

type store interface {
	DeleteItem(ctx context.Context, guid uuid.UUID) error
}
