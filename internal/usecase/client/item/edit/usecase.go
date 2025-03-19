// Package edit отвечает за сценарий изменения элемента на стороне клиента.
package edit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	dto "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

var ErrConflict = errors.New("conflict")

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

	newVersion, err := u.rpc.UpdateItem(ctx, &dto.UpdateItemIn{
		GUID:          item.GUID,
		EncryptedData: encryptedData,
		Version:       item.UpdatedAt.UTC().UnixMicro(),
	})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.FailedPrecondition {
			return ErrConflict
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	sec := newVersion / 1000000
	nsec := (newVersion % 1000000) * time.Microsecond.Nanoseconds()
	item.UpdatedAt = time.Unix(sec, nsec)

	err = u.itemStore.SaveItem(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
