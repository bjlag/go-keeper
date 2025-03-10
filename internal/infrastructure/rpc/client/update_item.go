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
}

func (c RPCClient) UpdateItem(ctx context.Context, in *UpdateItemIn) error {
	const op = "client.rpc.UpdateItem"

	rpcIn := &rpc.UpdateItemIn{
		Guid:          in.GUID.String(),
		EncryptedData: in.EncryptedData,
	}

	_, err := c.client.UpdateItem(ctx, rpcIn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
