package delete

import (
	"context"
	"errors"
	"fmt"

	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/item"
)

var ErrNotFoundUpdatedData = errors.New("deleted data is not found ")

type Usecase struct {
	store store
}

func NewUsecase(store store) *Usecase {
	return &Usecase{
		store: store,
	}
}

func (u Usecase) Do(ctx context.Context, in In) error {
	const op = "usecase.item.delete.Do"

	err := u.store.Delete(ctx, in.ItemGUID, in.UserGUID)
	if err != nil {
		if errors.Is(err, item.ErrNotAffectedRows) {
			return ErrNotFoundUpdatedData
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
