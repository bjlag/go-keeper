// Package backup отвечает за сброс локальных данных в бекап в зашифрованном виде.
package backup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

type Usecase struct {
	items  items
	tokens tokens
	backup backup
	cipher cipher
}

func NewUsecase(items items, tokens tokens, backups backup, cipher cipher) *Usecase {
	return &Usecase{
		items:  items,
		tokens: tokens,
		backup: backups,
		cipher: cipher,
	}
}

func (u *Usecase) Do(ctx context.Context) error {
	const op = "usecase.backup.Do"

	itemModels, err := u.items.Items(ctx, 100, 0)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(itemModels) == 0 {
		return nil
	}

	key := u.tokens.MasterKey()
	backupItems := make([]client.Backup, 0, len(itemModels))
	for _, i := range itemModels {
		v := &client.BackupValue{
			Category:  i.Category,
			Title:     i.Title,
			Value:     i.Value,
			Notes:     i.Notes,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		}

		marshaledValue, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		encrypted, err := u.cipher.Encrypt(marshaledValue, key)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		backupItems = append(backupItems, client.Backup{
			GUID:  i.GUID,
			Value: encrypted,
		})
	}

	if err = u.backup.Save(ctx, backupItems); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = u.items.EraseItems(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
