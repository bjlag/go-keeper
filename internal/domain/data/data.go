package data

import (
	"time"

	"github.com/google/uuid"
)

type Category int

const (
	CategoryLogin Category = iota
	CategoryText
	CategoryFile
	CategoryBankCard
)

type Data struct {
	GUID          uuid.UUID
	UserGUID      uuid.UUID
	EncryptedData []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
