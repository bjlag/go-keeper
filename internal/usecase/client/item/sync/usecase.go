package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

const prefixOp = "usecase.item.sync."

var ErrUnknownCategory = errors.New("unknown category")

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

func (u *Usecase) Do(ctx context.Context, guid uuid.UUID) (*model.Item, error) {
	const op = prefixOp + "Do"

	key := u.keyStore.MasterKey()

	in := &rpc.GetByGUIDIn{
		GUID: guid,
	}
	item, err := u.client.GetByGUID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	decrypted, err := u.cipher.Decrypt(item.EncryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var data model.EncryptedData
	err = json.Unmarshal(decrypted, &data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var value interface{}
	if data.Value != nil {
		switch data.Category {
		case model.CategoryPassword:
			value = &model.Password{}
		case model.CategoryText:
			break
		case model.CategoryFile:
			value = &model.File{}
		case model.CategoryBankCard:
			value = &model.BankCard{}
		default:
			return nil, fmt.Errorf("%w: %d", ErrUnknownCategory, data.Category)
		}

		err = json.Unmarshal(*data.Value, &value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	itemModel := model.Item{
		GUID:      item.GUID,
		Category:  data.Category,
		Title:     data.Title,
		Value:     value,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	err = u.itemStore.SaveItem(ctx, itemModel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &itemModel, nil
}
