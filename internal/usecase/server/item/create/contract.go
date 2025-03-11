package create

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

type store interface {
	Create(ctx context.Context, item model.Item) error
}
