package delete

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type Usecase struct {
	server server
	store  store
}

func NewUsecase(server server, store store) *Usecase {
	return &Usecase{
		server: server,
		store:  store,
	}
}

func (u *Usecase) Do(ctx context.Context, guid uuid.UUID) error {
	const op = "usecase.item.delete.Do"

	err := u.server.DeleteItem(ctx, &rpc.DeleteItemIn{
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
