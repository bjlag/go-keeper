package user

import (
	"time"

	"github.com/google/uuid"

	model "github.com/bjlag/go-keeper/internal/domain/server/user"
)

type row struct {
	GUID         uuid.UUID `db:"guid"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func convertToRow(m *model.User) row {
	return row{
		GUID:         m.GUID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (r row) convertToModel() *model.User {
	return &model.User{
		GUID:         r.GUID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
