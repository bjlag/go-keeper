package create

import (
	"time"

	"github.com/google/uuid"
)

type In struct {
	ItemGUID      uuid.UUID
	UserGUID      uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
}
