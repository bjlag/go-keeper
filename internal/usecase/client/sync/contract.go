package sync

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type client interface {
	GetAllItems(ctx context.Context, in *rpc.GetAllItemsIn) (*rpc.GetAllItemsOut, error)
}

type itemStore interface {
	SaveItems(ctx context.Context, items []model.RawItem) error
}

type keyStore interface {
	MasterKey() []byte
}

type cipher interface {
	Decrypt(data, key []byte) ([]byte, error)
}
