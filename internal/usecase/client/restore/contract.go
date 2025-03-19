package restore

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type items interface {
	SaveItems(ctx context.Context, items []model.RawItem) error
}

type tokens interface {
	MasterKey() []byte
}

type backup interface {
	Get(ctx context.Context) ([]model.Backup, error)
	Erase(ctx context.Context) error
}

type cipher interface {
	Decrypt(encryptedData, key []byte) ([]byte, error)
}
