package create

import (
	"context"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type Usecase struct {
	store store
}

func NewUsecase(store store) *Usecase {
	return &Usecase{
		store: store,
	}
}

func (u *Usecase) Do(ctx context.Context, item model.Item) error {
	const op = "usecase.item.create.Do"

	err := u.store.CreateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
