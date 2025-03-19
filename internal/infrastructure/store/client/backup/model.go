package backup

import (
	"github.com/google/uuid"

	"github.com/bjlag/go-keeper/internal/domain/client"
)

type row struct {
	GUID  uuid.UUID `db:"guid"`
	Value []byte    `db:"value"`
}

func fromModel(item client.Backup) row {
	return row{
		GUID:  item.GUID,
		Value: item.Value,
	}
}

func (r row) toModel() client.Backup {
	return client.Backup{
		GUID:  r.GUID,
		Value: r.Value,
	}
}

func toModels(rows []row) []client.Backup {
	items := make([]client.Backup, len(rows))
	for i, item := range rows {
		items[i] = item.toModel()
	}
	return items
}
