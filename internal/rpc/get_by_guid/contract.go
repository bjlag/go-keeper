package get_by_guid

import (
	"context"

	"github.com/bjlag/go-keeper/internal/domain/server/data"
	"github.com/bjlag/go-keeper/internal/usecase/server/item/get_by_guid"
)

type usecase interface {
	Do(ctx context.Context, data get_by_guid.Data) (*data.Item, error)
}
