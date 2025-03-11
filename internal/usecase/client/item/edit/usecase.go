package edit

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
	const op = "usecase.item.edit.Do"

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

	err = u.rpc.UpdateItem(ctx, &dto.UpdateItemIn{
		GUID:          item.GUID,
		EncryptedData: encryptedData,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.itemStore.SaveItem(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
