package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

// CreateItemIn данные создаваемого элемента.
type CreateItemIn struct {
	GUID          uuid.UUID // GUID создаваемого элемента.
	EncryptedData []byte    // EncryptedData секретные данные элемента в зашифрованном виде.
	CreatedAt     time.Time // CreatedAt дата и время создания элемента.
}

// CreateItem метод для создания элемента.
func (c RPCClient) CreateItem(ctx context.Context, in *CreateItemIn) error {
	const op = "client.rpc.CreateItem"

	rpcIn := &rpc.CreateItemIn{
		Guid:          in.GUID.String(),
		EncryptedData: in.EncryptedData,
		CreatedAt:     timestamppb.New(in.CreatedAt),
	}

	_, err := c.client.CreateItem(ctx, rpcIn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
