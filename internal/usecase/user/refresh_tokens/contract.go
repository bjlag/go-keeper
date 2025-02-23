package refresh_tokens

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/user"
)

type userStore interface {
	GetByGUID(ctx context.Context, guid string) (*model.User, error)
}

type tokenGenerator interface {
	GetUserGUID(tokenString string) (string, error)
	GenerateTokens(guid string) (accessToken, refreshToken string, err error)
}
