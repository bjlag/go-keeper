package save

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type store interface {
	SaveItem(ctx context.Context, item model.Item) error
}
