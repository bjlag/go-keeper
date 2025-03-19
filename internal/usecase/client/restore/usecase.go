// Package restore отвечает за восстановление локальных данных из бекапа.
package restore

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

	backupItems, err := u.backup.Get(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(backupItems) == 0 {
		return nil
	}

	var v client.BackupValue

	key := u.tokens.MasterKey()
	itemModels := make([]client.RawItem, 0, len(backupItems))
	for _, bi := range backupItems {
		decrypted, err := u.cipher.Decrypt(bi.Value, key)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = json.Unmarshal(decrypted, &v)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		itemModels = append(itemModels, client.RawItem{
			GUID:      bi.GUID,
			Category:  v.Category,
			Title:     v.Title,
			Value:     v.Value,
			Notes:     v.Notes,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	if err = u.items.SaveItems(ctx, itemModels); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = u.backup.Erase(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
