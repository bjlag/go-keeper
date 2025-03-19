package user

import (
	"time"

	"github.com/google/uuid"
)

// User описывает пользователя.
type User struct {
	// GUID уникальный идентификатор пользователя.
	GUID uuid.UUID
	// Email пользователя.
	Email string
	// PasswordHash хеш пароля.
	PasswordHash string
	// CreatedAt дата и время создания пользователя.
	CreatedAt time.Time
	// UpdatedAt дата и время обновления пользователя.
	UpdatedAt time.Time
}
