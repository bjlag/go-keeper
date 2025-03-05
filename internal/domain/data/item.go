package data

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	GUID          uuid.UUID
	UserGUID      uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
