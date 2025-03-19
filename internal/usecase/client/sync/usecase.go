// Package sync отвечает за сценарий синхронизации клиента с сервером.
package sync

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

const (
	prefixOp = "usecase.sync."

	limit = 40
)

type Usecase struct {
	client    client
	itemStore itemStore
	keyStore  keyStore
	cipher    cipher
}

func NewUsecase(client client, itemStore itemStore, keyStore keyStore, cipher cipher) *Usecase {
	return &Usecase{
		client:    client,
		itemStore: itemStore,
		keyStore:  keyStore,
		cipher:    cipher,
	}
}

func (u *Usecase) Do(ctx context.Context) error {
	const op = prefixOp + "Do"

	var offset uint32

	key := u.keyStore.MasterKey()

	for {
		in := &rpc.GetAllItemsIn{
			Limit:  limit,
			Offset: offset,
		}
		out, err := u.client.GetAllItems(ctx, in)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if out == nil || len(out.Items) == 0 {
			break
		}

		items := make([]model.RawItem, 0, len(out.Items))
		for _, item := range out.Items {
			decrypted, err := u.cipher.Decrypt(item.EncryptedData, key)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			var data model.EncryptedData
			err = json.Unmarshal(decrypted, &data)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			items = append(items, model.RawItem{
				GUID:      item.GUID,
				Category:  data.Category,
				Title:     data.Title,
				Value:     data.Value,
				Notes:     data.Notes,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		err = u.itemStore.SaveItems(ctx, items)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		offset += limit
	}

	return nil
}
