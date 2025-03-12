package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type UpdateItemIn struct {
	GUID          uuid.UUID
	EncryptedData []byte
	Version       int64
}

func (c RPCClient) UpdateItem(ctx context.Context, in *UpdateItemIn) (int64, error) {
	const op = "client.rpc.UpdateItem"

	rpcIn := &rpc.UpdateItemIn{
		Guid:          in.GUID.String(),
		EncryptedData: in.EncryptedData,
		Version:       in.Version,
	}

	rpcOut, err := c.client.UpdateItem(ctx, rpcIn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return rpcOut.GetNewVersion(), nil
}
