package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

// UpdateItemIn параметры запроса.
type UpdateItemIn struct {
	GUID          uuid.UUID // GUID идентификатор обновляемого элемента.
	EncryptedData []byte    // EncryptedData зашифрованные данные элемента
	Version       int64     // Version версия, с которой обновляем элемент (текущая версия).
}

// UpdateItem метод для обновления элемента.
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
