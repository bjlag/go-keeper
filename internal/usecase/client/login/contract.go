package login

import (
	"context"

	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type client interface {
	Login(ctx context.Context, in rpc.LoginIn) (*rpc.LoginOut, error)
}

type tokens interface {
	SaveTokens(accessToken, refreshToken string)
}
