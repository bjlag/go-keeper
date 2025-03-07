package create

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type store interface {
	CreateItem(ctx context.Context, item model.Item) error
}
