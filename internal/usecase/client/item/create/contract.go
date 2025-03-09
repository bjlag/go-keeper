package create

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type server interface {
	CreateItem(ctx context.Context, in *rpc.CreateItemIn) error
}

type store interface {
	CreateItem(ctx context.Context, item model.Item) error
}
