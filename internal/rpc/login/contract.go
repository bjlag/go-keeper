package login

import (
	"context"

	"github.com/bjlag/go-keeper/internal/usecase/user/login"
)

type usecase interface {
	Do(ctx context.Context, data login.Data) (*login.Result, error)
}
