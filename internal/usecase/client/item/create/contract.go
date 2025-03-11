package create

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	dto "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type rpc interface {
	CreateItem(ctx context.Context, in *dto.CreateItemIn) error
}

type itemStore interface {
	CreateItem(ctx context.Context, item model.Item) error
}

type keyStore interface {
	MasterKey() []byte
}

type cipher interface {
	Encrypt(data, key []byte) ([]byte, error)
}
