package get_all_data

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/data/get_all"
)

type usecase interface {
	Do(ctx context.Context, data get_all.Data) (*get_all.Result, error)
}
