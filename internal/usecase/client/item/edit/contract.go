package edit

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	dto "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type rpc interface {
	UpdateItem(ctx context.Context, in *dto.UpdateItemIn) (int64, error)
}

type itemStore interface {
	SaveItem(ctx context.Context, item model.Item) error
}

type keyStore interface {
	MasterKey() []byte
}

type cipher interface {
	Encrypt(data, key []byte) ([]byte, error)
}
