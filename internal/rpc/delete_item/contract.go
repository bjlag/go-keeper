package delete_item

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/item/remove"
)

type usecase interface {
	Do(ctx context.Context, data remove.In) error
}
