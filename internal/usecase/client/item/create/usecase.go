package create

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
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

func (u *Usecase) Do(ctx context.Context, item model.Item) error {
	const op = "usecase.item.create.Do"

	var value *[]byte
	if item.Value != nil {
		v, err := json.Marshal(item.Value)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		value = &v
	}

	data := model.EncryptedData{
		Title:    item.Title,
		Category: item.Category,
		Value:    value,
		Notes:    item.Notes,
	}

	encrypted, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.server.CreateItem(ctx, &rpc.CreateItemIn{
		GUID:          item.GUID,
		EncryptedData: encrypted,
		CreatedAt:     item.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.store.CreateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
