package client

import (
	"context"
	"fmt"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type LoginIn struct {
	Email    string
	Password string
}

type LoginOut struct {
	AccessToken  string
	RefreshToken string
}

func (c RPCClient) Login(ctx context.Context, in LoginIn) (*LoginOut, error) {
	const op = "client.rpc.Login"

	rpcIn := &rpc.LoginIn{
		Email:    in.Email,
		Password: in.Password,
	}

	out, err := c.client.Login(ctx, rpcIn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &LoginOut{
		AccessToken:  out.GetAccessToken(),
		RefreshToken: out.GetRefreshToken(),
	}, nil
}
