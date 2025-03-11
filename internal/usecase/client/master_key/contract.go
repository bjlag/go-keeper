package master_key

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	"github.com/bjlag/go-keeper/internal/infrastructure/crypt/master_key"
)

type tokens interface {
	SaveMasterKey(key []byte)
}

type options interface {
	OptionBySlug(ctx context.Context, slug string) (*model.Option, error)
	SaveOption(ctx context.Context, option model.Option) error
}

type salter interface {
	GenerateSalt() (*master_key.Salt, error)
}

type keymaker interface {
	GenerateMasterKey(password, salt []byte) []byte
}
