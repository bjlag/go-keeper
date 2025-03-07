package create

import (
	"github.com/google/uuid"
)

type In struct {
	ItemGUID      uuid.UUID
	UserGUID      uuid.UUID
	EncryptedData []byte
}
