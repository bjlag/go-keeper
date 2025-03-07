package create

import (
	"context"
	"fmt"
	"time"

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

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	err := u.store.CreateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
