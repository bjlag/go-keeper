package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

// GetByGUIDIn параметры запроса.
type GetByGUIDIn struct {
	GUID uuid.UUID // GUID элемента, данные которого надо получить.
}

// GetByGUIDOut результат работы.
type GetByGUIDOut struct {
	GUID          uuid.UUID // GUID идентификатор элемента.
	EncryptedData []byte    // EncryptedData зашифрованные данные элемента.
	CreatedAt     time.Time // CreatedAt дата и время создания элемента.
	UpdatedAt     time.Time // UpdatedAt дата и время обновления элемента.
}

// GetByGUID метод для получения элемента по его GUID.
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
