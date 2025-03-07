package create_item

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/data/create"
)

type usecase interface {
	Do(ctx context.Context, data create.In) error
}
