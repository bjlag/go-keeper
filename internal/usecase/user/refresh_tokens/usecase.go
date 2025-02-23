package refresh_tokens

import (
	"context"
	"errors"
	"fmt"

	"github.com/bjlag/go-keeper/internal/infrastructure/auth/jwt"
	storeUser "github.com/bjlag/go-keeper/internal/infrastructure/store/user"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type Usecase struct {
	userStore      userStore
	tokenGenerator tokenGenerator
}

func NewUsecase(userStore userStore, tokenGenerator tokenGenerator) *Usecase {
	return &Usecase{
		userStore:      userStore,
		tokenGenerator: tokenGenerator,
	}
}

func (u Usecase) Do(ctx context.Context, data Data) (*Result, error) {
	const op = "usecase.refreshTokens.Do"

	guid, err := u.tokenGenerator.GetUserGUID(data.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := u.userStore.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, storeUser.ErrNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, refreshToken, err := u.tokenGenerator.GenerateTokens(user.GUID.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Result{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
