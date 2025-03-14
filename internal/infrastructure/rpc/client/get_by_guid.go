package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type GetByGUIDIn struct {
	GUID uuid.UUID
}

type GetByGUIDOut struct {
	GUID          uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c RPCClient) GetByGUID(ctx context.Context, in *GetByGUIDIn) (*GetByGUIDOut, error) {
	const op = "client.rpc.GetByGUID"

	rpcIn := &rpc.GetByGuidIn{
		Guid: in.GUID.String(),
	}

	out, err := c.client.GetByGuid(ctx, rpcIn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	itemGUID, err := uuid.Parse(out.GetGuid())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &GetByGUIDOut{
		GUID:          itemGUID,
		EncryptedData: out.GetEncryptedData(),
		CreatedAt:     out.GetCreatedAt().AsTime(),
		UpdatedAt:     out.GetUpdatedAt().AsTime(),
	}, nil
}
