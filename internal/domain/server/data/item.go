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

type UpdatedItem struct {
	EncryptedData []byte
	UpdatedAt     time.Time
}
