package sync

import (
	"context"
	"encoding/json"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
)

const (
	prefixOp = "usecase.sync."

	limit = 1
)

type Data struct {
	Title      string         `json:"title"`
	CategoryID model.Category `json:"category_id"`
	Value      *[]byte        `json:"value,omitempty"`
	Notes      string         `json:"notes"`
}

func (d *Data) UnmarshalJSON(data []byte) error {
	type Alias Data

	alias := &struct {
		*Alias
		Value *json.RawMessage `json:"value,omitempty"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, alias); err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	if alias.Value != nil {
		value := []byte(*alias.Value)

		if alias.CategoryID == model.CategoryPassword {

		}
		d.Value = &value
	}

	return nil
}

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

		items := make([]model.RawItem, 0, len(out.Items))
		for _, item := range out.Items {
			// todo расшифровка
			// todo общие данные в отдельных полях

			var data Data
			err = json.Unmarshal(item.EncryptedData, &data)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			items = append(items, model.RawItem{
				GUID:       item.GUID,
				CategoryID: data.CategoryID,
				Title:      data.Title,
				Value:      data.Value,
				Notes:      data.Notes,
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  item.UpdatedAt,
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
