package get_all_items

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/item/get_all"
)

type usecase interface {
	Do(ctx context.Context, data get_all.Data) (*get_all.Result, error)
}
