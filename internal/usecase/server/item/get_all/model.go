package get_all

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	UserGUID uuid.UUID
	Limit    uint32
	Offset   uint32
}

type Result struct {
	Items []Item
}

type Item struct {
	GUID          uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
