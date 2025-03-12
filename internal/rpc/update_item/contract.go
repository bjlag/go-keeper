package update_item

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/item/update"
)

type usecase interface {
	Do(ctx context.Context, data update.In) (newVersion int64, err error)
}
