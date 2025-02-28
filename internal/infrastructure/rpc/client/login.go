package client

import (
	"context"

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
	rpcIn := &rpc.LoginIn{
		Email:    in.Email,
		Password: in.Password,
	}

	out, err := c.client.Login(ctx, rpcIn)
	if err != nil {
		return nil, err
	}

	return &LoginOut{
		AccessToken:  out.GetAccessToken(),
		RefreshToken: out.GetRefreshToken(),
	}, nil
}
