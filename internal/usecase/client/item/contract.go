package item

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type itemStore interface {
	ItemsByCategory(ctx context.Context, category model.Category) ([]model.RawItem, error)
}
