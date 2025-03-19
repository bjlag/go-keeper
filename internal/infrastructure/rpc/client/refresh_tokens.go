package client

import (
	"context"
	"fmt"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

// RefreshTokensIn параметры запроса.
type RefreshTokensIn struct {
	RefreshToken string // RefreshToken refresh токен.
}

// RefreshTokensOut результат.
type RefreshTokensOut struct {
	AccessToken  string // AccessToken access токен.
	RefreshToken string // RefreshToken refresh токен.
}

// RefreshTokens метод для обновления токенов используя refresh токен.
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
