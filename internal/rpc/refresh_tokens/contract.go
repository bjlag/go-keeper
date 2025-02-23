package refresh_tokens

import (
	"context"

	rt "github.com/bjlag/go-keeper/internal/usecase/user/refresh_tokens"
)

type usecase interface {
	Do(ctx context.Context, data rt.Data) (*rt.Result, error)
}
