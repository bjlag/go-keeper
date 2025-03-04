package client

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

type Item struct {
	GUID      uuid.UUID
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
