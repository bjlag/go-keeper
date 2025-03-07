package update_item

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/data/update"
)

type usecase interface {
	Do(ctx context.Context, data update.In) error
}
