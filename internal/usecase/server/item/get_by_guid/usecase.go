// Package get_by_guid отвечает за получение элемента на стороне сервера.
package get_by_guid

import (
	"context"
	"errors"
	"fmt"

	"github.com/bjlag/go-keeper/internal/domain/server/data"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/item"
)

var ErrNoData = errors.New("no data")

type Usecase struct {
	dataStore dataStore
}

func NewUsecase(dataStore dataStore) *Usecase {
	return &Usecase{
		dataStore: dataStore,
	}
}

func (u Usecase) Do(ctx context.Context, data Data) (*data.Item, error) {
	const op = "usecase.item.getByGuid.Do"

	model, err := u.dataStore.UserItemByGUID(ctx, data.UserGUID, data.GUID)
	if err != nil {
		if errors.Is(err, item.ErrNotFound) {
			return nil, ErrNoData
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return model, nil
}
