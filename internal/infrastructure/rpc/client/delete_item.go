package client

import (
	"context"
	"fmt"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/google/uuid"
)

type DeleteItemIn struct {
	GUID uuid.UUID
}

func (c RPCClient) DeleteItem(ctx context.Context, in *DeleteItemIn) error {
	const op = "client.rpc.DeleteItem"

	rpcIn := &rpc.DeleteItemIn{
		Guid: in.GUID.String(),
	}

	_, err := c.client.DeleteItem(ctx, rpcIn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
