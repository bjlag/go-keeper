package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/bjlag/go-keeper/internal/domain/data"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/item"
)

var ErrNotFoundUpdatedData = errors.New("not found updated data")

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

	data := model.UpdatedItem{
		EncryptedData: in.EncryptedData,
		UpdatedAt:     time.Now(),
	}

	err := u.store.Update(ctx, in.ItemGUID, in.UserGUID, data)
	if err != nil {
		if errors.Is(err, item.ErrNotAffectedRows) {
			return ErrNotFoundUpdatedData
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
