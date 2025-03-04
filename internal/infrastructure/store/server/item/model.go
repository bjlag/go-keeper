package item

import (
	"time"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/domain/data"
)

type Row struct {
	GUID          uuid.UUID `db:"guid"`
	UserGUID      uuid.UUID `db:"user_guid"`
	EncryptedData []byte    `db:"encrypted_data"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (r *Row) convertToModel() data.Item {
	return data.Item{
		GUID:          r.GUID,
		UserGUID:      r.UserGUID,
		EncryptedData: r.EncryptedData,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func convertToModels(rows []Row) []data.Item {
	result := make([]data.Item, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.convertToModel())
	}
	return result
}
