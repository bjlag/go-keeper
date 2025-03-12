package get_by_guid

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	GUID     uuid.UUID
	UserGUID uuid.UUID
}

type Result struct {
	GUID          uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
