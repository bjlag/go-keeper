package item

import (
	"time"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type row struct {
	GUID       uuid.UUID      `db:"guid"`
	CategoryID model.Category `db:"category_id"`
	Title      string         `db:"title"`
	Value      *[]byte        `db:"value"`
	Notes      string         `db:"notes"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}

func toModels(rows []row) []model.RawItem {
	items := make([]model.RawItem, len(rows))
	for i, r := range rows {
		items[i] = r.toModel()
	}
	return items
}

func (r *row) toModel() model.RawItem {
	return model.RawItem{
		GUID:       r.GUID,
		CategoryID: r.CategoryID,
		Title:      r.Title,
		Value:      r.Value,
		Notes:      r.Notes,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}
