package sync

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type client interface {
	GetAllItems(ctx context.Context, in *rpc.GetAllItemsIn) (*rpc.GetAllItemsOut, error)
}

type store interface {
	SaveItems(ctx context.Context, items []model.Item) error
}
