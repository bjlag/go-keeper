package login

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	storeUser "github.com/bjlag/go-keeper/internal/infrastructure/store/user"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPasswordIncorrect = errors.New("password incorrect")
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
	const op = "usecase.login.Do"

	user, err := u.userStore.GetByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, storeUser.ErrNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password))
	if err != nil {
		return nil, ErrPasswordIncorrect
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
