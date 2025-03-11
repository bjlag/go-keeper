package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	GUID         uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
