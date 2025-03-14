// Package create отвечает за сценарий создания элемента на стороне сервера.
package create

import (
	"context"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

type Usecase struct {
	store store
}

func NewUsecase(store store) *Usecase {
	return &Usecase{
		store: store,
	}
}

func (u Usecase) Do(ctx context.Context, in In) error {
	const op = "usecase.item.create.Do"

	data := model.Item{
		GUID:          in.ItemGUID,
		UserGUID:      in.UserGUID,
		EncryptedData: in.EncryptedData,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.CreatedAt,
	}

	err := u.store.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
