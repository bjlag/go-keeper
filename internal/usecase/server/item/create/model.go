package create

import (
	"github.com/google/uuid"
	"time"
)

type In struct {
	ItemGUID      uuid.UUID
	UserGUID      uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
}
