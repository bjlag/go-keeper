package login

import (
	"context"
	"fmt"

	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type Usecase struct {
	client client
}

func NewUsecase(client client) *Usecase {
	return &Usecase{
		client: client,
	}
}

func (u *Usecase) Do(ctx context.Context, data Data) (*Result, error) {
	const op = "usecase.login.Do"

	out, err := u.client.Login(ctx, rpc.LoginIn{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Result{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, nil
}
