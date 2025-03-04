package sync

import (
	"context"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

const (
	prefixOp = "usecase.sync."

	limit = 1
)

type Usecase struct {
	client client
	store  store
}

func NewUsecase(client client, store store) *Usecase {
	return &Usecase{
		client: client,
		store:  store,
	}
}

func (u *Usecase) Do(ctx context.Context) error {
	const op = prefixOp + "Do"

	var offset uint32

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

		items := make([]model.Item, 0, len(out.Items))
		for _, item := range out.Items {
			// todo расшифровка
			// todo общие данные в отдельных полях
			items = append(items, model.Item{
				GUID:      item.GUID,
				Data:      item.EncryptedData,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		err = u.store.SaveItems(ctx, items)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		offset += limit
	}

	return nil
}
