package delete

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Usecase struct {
	store store
}

func NewUsecase(store store) *Usecase {
	return &Usecase{
		store: store,
	}
}

func (u *Usecase) Do(ctx context.Context, guid uuid.UUID) error {
	const op = "usecase.item.delete.Do"

	err := u.store.DeleteItem(ctx, guid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
