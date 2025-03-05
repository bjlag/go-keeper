package client

import (
	"time"

	"github.com/google/uuid"
)

type Category int

const (
	CategoryLogin Category = iota
	CategoryText
	CategoryBlob
	CategoryBankCard
)

type Item struct {
	GUID       uuid.UUID
	CategoryID Category
	Title      string
	Value      *[]byte
	Notes      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
