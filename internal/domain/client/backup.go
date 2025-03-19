package client

import (
	"github.com/google/uuid"
	"time"
)

type Backup struct {
	GUID  uuid.UUID
	Value []byte
}

type BackupValue struct {
	Category  Category  `json:"category"`
	Title     string    `json:"title"`
	Value     *[]byte   `json:"value"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
