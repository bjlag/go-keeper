package item

import (
	"time"

	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/domain/data"
)

type row struct {
	GUID          uuid.UUID `db:"guid"`
	UserGUID      uuid.UUID `db:"user_guid"`
	EncryptedData []byte    `db:"encrypted_data"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (r *row) convertToModel() data.Item {
	return data.Item{
		GUID:          r.GUID,
		UserGUID:      r.UserGUID,
		EncryptedData: r.EncryptedData,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func convertToModels(rows []row) []data.Item {
	result := make([]data.Item, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.convertToModel())
	}
	return result
}

type updated struct {
	GUID          uuid.UUID `db:"guid"`
	UserGUID      uuid.UUID `db:"user_guid"`
	EncryptedData []byte    `db:"encrypted_data"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type deleted struct {
	GUID     uuid.UUID `db:"guid"`
	UserGUID uuid.UUID `db:"user_guid"`
}
