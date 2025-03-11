package register

import (
	"context"
	"fmt"

	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type Usecase struct {
	client client
	tokens tokens
}

func NewUsecase(client client, tokens tokens) *Usecase {
	return &Usecase{
		client: client,
		tokens: tokens,
	}
}

func (u *Usecase) Do(ctx context.Context, data Data) error {
	const op = "usecase.register.Do"

	out, err := u.client.Register(ctx, rpc.RegisterIn{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	u.tokens.SaveTokens(out.AccessToken, out.RefreshToken)

	return nil
}
