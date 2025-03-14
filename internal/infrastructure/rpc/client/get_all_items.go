package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

// GetAllItemsIn параметры запроса.
type GetAllItemsIn struct {
	Limit  uint32 // Limit сколько получить данных максимум.
	Offset uint32 // Offset с какой позиции получать данные.
}

// GetAllItemsOut результат работы.
type GetAllItemsOut struct {
	Items []GetAllDataItem // Items список полученных элементов.
}

// GetAllDataItem данные полученного элемента.
type GetAllDataItem struct {
	GUID          uuid.UUID // GUID идентификатор.
	EncryptedData []byte    // EncryptedData зашифрованные данные.
	CreatedAt     time.Time // CreatedAt дата и время создания.
	UpdatedAt     time.Time // UpdatedAt дата и время обновления.
}

// GetAllItems метод для получения всех данных авторизованного пользователя.
func (c RPCClient) GetAllItems(ctx context.Context, in *GetAllItemsIn) (*GetAllItemsOut, error) {
	const op = "client.rpc.GetAllItems"

	rpcIn := &rpc.GetAllItemsIn{
		Limit:  in.Limit,
		Offset: in.Offset,
	}

	out, err := c.client.GetAllItems(ctx, rpcIn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if out == nil || len(out.GetItems()) == 0 {
		return nil, nil
	}

	items := make([]GetAllDataItem, len(out.GetItems()))
	for i, item := range out.GetItems() {
		guid, err := uuid.Parse(item.GetGuid())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		items[i] = GetAllDataItem{
			GUID:          guid,
			EncryptedData: item.GetEncryptedData(),
			CreatedAt:     item.GetCreatedAt().AsTime(),
			UpdatedAt:     item.GetUpdatedAt().AsTime(),
		}
	}

	return &GetAllItemsOut{
		Items: items,
	}, nil
}
