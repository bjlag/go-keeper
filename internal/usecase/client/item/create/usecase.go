// Package create отвечает за сценарий создания элемента на стороне клиента.
package create

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	dto "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

type Usecase struct {
	rpc       rpc
	itemStore itemStore
	keyStore  keyStore
	cipher    cipher
}

func NewUsecase(rpc rpc, itemStore itemStore, keyStore keyStore, cipher cipher) *Usecase {
	return &Usecase{
		rpc:       rpc,
		itemStore: itemStore,
		keyStore:  keyStore,
		cipher:    cipher,
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

	plainText, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	encryptedData, err := u.cipher.Encrypt(plainText, u.keyStore.MasterKey())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.rpc.CreateItem(ctx, &dto.CreateItemIn{
		GUID:          item.GUID,
		EncryptedData: encryptedData,
		CreatedAt:     item.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.itemStore.CreateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
