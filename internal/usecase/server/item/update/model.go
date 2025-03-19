package update

import (
	"github.com/google/uuid"
)

type In struct {
	UserGUID      uuid.UUID
	ItemGUID      uuid.UUID
	EncryptedData []byte
	Version       int64
}
