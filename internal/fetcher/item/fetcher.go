// Package item отвечает за получение данных по элементам.
package item

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

const prefixOp = "fetcher.item"

// ErrUnknownCategory ошибка на случай если получены данные, значение которых не известны.
// Неизвестно, в какую модель их раскладывать.
var ErrUnknownCategory = errors.New("unknown category")

type Fetcher struct {
	itemStore itemStore
}

func NewFetcher(itemStore itemStore) *Fetcher {
	return &Fetcher{
		itemStore: itemStore,
	}
}

// ItemsByCategory получает элементы из локальной базы клиента по указанной категории.
func (u *Fetcher) ItemsByCategory(ctx context.Context, category model.Category) ([]model.Item, error) {
	const op = prefixOp + ".ItemsByCategory"

	rawItems, err := u.itemStore.ItemsByCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	items := make([]model.Item, len(rawItems))
	for i, item := range rawItems {
		var v interface{}
		if item.Value != nil {
			switch item.Category {
			case model.CategoryPassword:
				v = &model.Password{}
			case model.CategoryText:
				break
			case model.CategoryFile:
				v = &model.File{}
			case model.CategoryBankCard:
				v = &model.BankCard{}
			default:
				return nil, fmt.Errorf("%w: %d", ErrUnknownCategory, item.Category)
			}

			err = json.Unmarshal(*item.Value, &v)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}

		items[i] = model.Item{
			GUID:      item.GUID,
			Category:  item.Category,
			Title:     item.Title,
			Value:     v,
			Notes:     item.Notes,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return items, nil
}
