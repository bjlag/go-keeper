package client

import (
	"context"
	"fmt"
	"time"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type GetAllItemsIn struct {
	Limit  uint32
	Offset uint32
}

type GetAllItemsOut struct {
	Items []GetAllDataItem
}

type GetAllDataItem struct {
	GUID          string
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

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
		items[i] = GetAllDataItem{
			GUID:          item.GetGuid(),
			EncryptedData: item.GetEncryptedData(),
			CreatedAt:     item.GetCreatedAt().AsTime(),
			UpdatedAt:     item.GetUpdatedAt().AsTime(),
		}
	}

	return &GetAllItemsOut{
		Items: items,
	}, nil
}
