package refresh_tokens

import (
	"context"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/server/user"
)

type userStore interface {
	GetByGUID(ctx context.Context, guid uuid.UUID) (*model.User, error)
}

type tokenGenerator interface {
	GetUserGUIDFromRefreshToken(tokenString string) (uuid.UUID, error)
	GenerateTokens(guid string) (accessToken, refreshToken string, err error)
}
