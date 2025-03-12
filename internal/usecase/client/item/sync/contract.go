package sync

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type client interface {
	GetByGUID(ctx context.Context, in *rpc.GetByGUIDIn) (*rpc.GetByGUIDOut, error)
}

type itemStore interface {
	SaveItem(ctx context.Context, item model.Item) error
}

type keyStore interface {
	MasterKey() []byte
}

type cipher interface {
	Decrypt(data, key []byte) ([]byte, error)
}
