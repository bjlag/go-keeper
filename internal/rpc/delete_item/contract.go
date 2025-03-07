package delete_item

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/item/delete"
)

type usecase interface {
	Do(ctx context.Context, data delete.In) error
}
