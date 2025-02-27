package register

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	model "github.com/bjlag/go-keeper/internal/domain/user"
	storeUser "github.com/bjlag/go-keeper/internal/infrastructure/store/user"
)

var ErrUserAlreadyExists = errors.New("user already exists")

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
	const op = "usecase.register.Do"

	exist, err := u.isUserExists(ctx, data.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if exist {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user := &model.User{
		GUID:         uuid.New(),
		Email:        data.Email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = u.userStore.Add(ctx, user)
	if err != nil {
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

func (u Usecase) isUserExists(ctx context.Context, email string) (bool, error) {
	m, err := u.userStore.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storeUser.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user by email: %w", err)
	}

	return m != nil, nil
}
