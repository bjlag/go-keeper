package client

import (
	"context"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type RegisterIn struct {
	Email    string
	Password string
}

type RegisterOut struct {
	AccessToken  string
	RefreshToken string
}

func (c RPCClient) Register(ctx context.Context, in RegisterIn) (*RegisterOut, error) {
	rpcIn := &rpc.RegisterIn{
		Email:    in.Email,
		Password: in.Password,
	}

	out, err := c.client.Register(ctx, rpcIn)
	if err != nil {
		return nil, err
	}

	return &RegisterOut{
		AccessToken:  out.GetAccessToken(),
		RefreshToken: out.GetRefreshToken(),
	}, nil
}
