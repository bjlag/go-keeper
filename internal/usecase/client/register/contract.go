package register

import (
	"context"

	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type client interface {
	Register(ctx context.Context, in rpc.RegisterIn) (*rpc.RegisterOut, error)
}
