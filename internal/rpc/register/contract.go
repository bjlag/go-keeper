package register

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/server/user/register"
)

type usecase interface {
	Do(ctx context.Context, data register.Data) (*register.Result, error)
}
