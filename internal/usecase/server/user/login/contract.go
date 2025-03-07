package login

import (
	"context"

	model "github.com/bjlag/go-keeper/internal/domain/server/user"
)

type userStore interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type tokenGenerator interface {
	GenerateTokens(guid string) (accessToken, refreshToken string, err error)
}
