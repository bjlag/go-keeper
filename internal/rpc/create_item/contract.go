package create_item

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/item/create"
)

type usecase interface {
	Do(ctx context.Context, data create.In) error
}
