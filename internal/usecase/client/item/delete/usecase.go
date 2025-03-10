package delete

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	dto "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type Usecase struct {
	rpc   rpc
	store store
}

func NewUsecase(rpc rpc, store store) *Usecase {
	return &Usecase{
		rpc:   rpc,
		store: store,
	}
}

func (u *Usecase) Do(ctx context.Context, guid uuid.UUID) error {
	const op = "usecase.item.delete.Do"

	err := u.rpc.DeleteItem(ctx, &dto.DeleteItemIn{
		GUID: guid,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.store.DeleteItem(ctx, guid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
