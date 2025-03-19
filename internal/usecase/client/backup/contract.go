package backup

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type items interface {
	Items(ctx context.Context, limit, offset int64) ([]model.RawItem, error)
	EraseItems(ctx context.Context) error
}

type tokens interface {
	MasterKey() []byte
}

type backup interface {
	Save(ctx context.Context, items []model.Backup) error
}

type cipher interface {
	Encrypt(data, key []byte) ([]byte, error)
}
