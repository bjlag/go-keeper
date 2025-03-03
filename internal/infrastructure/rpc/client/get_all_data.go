package client

import (
	"context"
	"fmt"
	"time"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type GetAllDataIn struct {
	Limit  uint32
	Offset uint32
}

type GetAllDataOut struct {
	Items []GetAllDataItem
}

type GetAllDataItem struct {
	GUID          string
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c RPCClient) GetAllData(ctx context.Context, in GetAllDataIn) (*GetAllDataOut, error) {
	const op = "client.rpc.GetAllData"

	rpcIn := &rpc.GetAllDataIn{
		Limit:  in.Limit,
		Offset: in.Offset,
	}

	out, err := c.client.GetAllData(ctx, rpcIn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if out == nil || len(out.Items) == 0 {
		return nil, nil
	}

	items := make([]GetAllDataItem, len(out.Items))
	for i, item := range out.Items {
		items[i] = GetAllDataItem{
			GUID:          item.Guid,
			EncryptedData: item.EncryptedData,
			CreatedAt:     item.CreatedAt.AsTime(),
			UpdatedAt:     item.UpdatedAt.AsTime(),
		}
	}

	return &GetAllDataOut{
		Items: items,
	}, nil
}
