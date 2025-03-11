package remove

import (
	"github.com/google/uuid"
)

type In struct {
	UserGUID uuid.UUID
	ItemGUID uuid.UUID
}
