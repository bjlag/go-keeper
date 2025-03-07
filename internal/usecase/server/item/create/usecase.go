package create

import (
	"context"
	"fmt"
	"time"

	model "github.com/bjlag/go-keeper/internal/domain/server/data"
)

type Usecase struct {
	store store
}

func NewUsecase(store store) *Usecase {
	return &Usecase{
		store: store,
	}
}

func (u Usecase) Do(ctx context.Context, in In) error {
	const op = "usecase.item.update.Do"

	data := model.Item{
		GUID:          in.ItemGUID,
		UserGUID:      in.UserGUID,
		EncryptedData: in.EncryptedData,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := u.store.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
