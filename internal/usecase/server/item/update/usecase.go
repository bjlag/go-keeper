// Package update отвечает за обновление элемента на стороне сервера.
package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/item"
)

var (
	ErrNotFoundUpdatedData = errors.New("updated data is not found")
	ErrConflict            = errors.New("conflict")
)

type Usecase struct {
	store store
}

func NewUsecase(store store) *Usecase {
	return &Usecase{
		store: store,
	}
}

func (u Usecase) Do(ctx context.Context, in In) (newVersion int64, err error) {
	const op = "usecase.item.update.Do"

	currentItem, err := u.store.ItemByGUID(ctx, in.ItemGUID)
	if err != nil {
		if errors.Is(err, item.ErrNotFound) {
			return 0, ErrNotFoundUpdatedData
		}
	}

	if in.Version != currentItem.UpdatedAt.UTC().UnixMicro() {
		return 0, ErrConflict
	}

	updatedAt := time.Now()
	data := model.UpdatedItem{
		EncryptedData: in.EncryptedData,
		UpdatedAt:     updatedAt,
	}

	err = u.store.Update(ctx, in.ItemGUID, in.UserGUID, data)
	if err != nil {
		if errors.Is(err, item.ErrNotAffectedRows) {
			return 0, ErrNotFoundUpdatedData
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return updatedAt.UnixMicro(), nil
}
