package client

import (
	"context"
	"fmt"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

// RegisterIn параметры запроса.
type RegisterIn struct {
	Email    string // Email пользователя.
	Password string // Password пароль пользователя.
}

// RegisterOut результат.
type RegisterOut struct {
	AccessToken  string // AccessToken access токен.
	RefreshToken string // RefreshToken refresh токен.
}

// Register метод для регистрации пользователя.
func (c RPCClient) Register(ctx context.Context, in RegisterIn) (*RegisterOut, error) {
	const op = "client.rpc.Register"

	rpcIn := &rpc.RegisterIn{
		Email:    in.Email,
		Password: in.Password,
	}

	out, err := c.client.Register(ctx, rpcIn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RegisterOut{
		AccessToken:  out.GetAccessToken(),
		RefreshToken: out.GetRefreshToken(),
	}, nil
}
