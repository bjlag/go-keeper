package delete

import (
	"context"

	"github.com/google/uuid"

	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type server interface {
	DeleteItem(ctx context.Context, in *rpc.DeleteItemIn) error
}

type store interface {
	DeleteItem(ctx context.Context, guid uuid.UUID) error
}
