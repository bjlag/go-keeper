package client

import (
	"context"
	"fmt"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type RefreshTokensIn struct {
	RefreshToken string
}

type RefreshTokensOut struct {
	AccessToken  string
	RefreshToken string
}

func (c RPCClient) RefreshTokens(ctx context.Context, in RefreshTokensIn) (*RefreshTokensOut, error) {
	const op = "client.rpc.RefreshTokens"

	rpcIn := &rpc.RefreshTokensIn{
		RefreshToken: in.RefreshToken,
	}

	out, err := c.client.RefreshTokens(ctx, rpcIn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RefreshTokensOut{
		AccessToken:  out.GetAccessToken(),
		RefreshToken: out.GetRefreshToken(),
	}, nil
}
