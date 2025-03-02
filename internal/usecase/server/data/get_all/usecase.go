package get_all

import (
	"context"
	"errors"
	"fmt"

	storeUser "github.com/bjlag/go-keeper/internal/infrastructure/store/user"
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

func (u Usecase) Do(ctx context.Context, data Data) (*Result, error) {
	const op = "usecase.data.getAll.Do"

	rows, err := u.dataStore.GetAllByUser(ctx, data.UserGUID, data.Limit, data.Offset)
	if err != nil {
		if errors.Is(err, storeUser.ErrNotFound) {
			return nil, ErrNoData
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	items := make([]Item, 0, len(rows))
	for _, r := range rows {
		items = append(items, Item{
			GUID:          r.GUID,
			EncryptedData: r.EncryptedData,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		})
	}

	return &Result{
		Items: items,
	}, nil
}
