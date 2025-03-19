package data

import (
	"time"

	"github.com/google/uuid"
)

// Item описывает секретные данные.
type Item struct {
	// GUID уникальный идентификатор.
	GUID uuid.UUID
	// UserGUID идентификатор пользователя, которому принадлежат данные.
	UserGUID uuid.UUID
	// EncryptedData сами секретные данные в зашифрованном виде.
	EncryptedData []byte
	// CreatedAt дата и время создания записи.
	CreatedAt time.Time
	// UpdatedAt дата и время обновления записи.
	UpdatedAt time.Time
}

// UpdatedItem описывает данные для обновления.
type UpdatedItem struct {
	EncryptedData []byte
	UpdatedAt     time.Time
}
