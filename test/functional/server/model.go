package server

import "time"

type Item struct {
	GUID          string    `db:"guid"`
	UserGUID      string    `db:"user_guid"`
	EncryptedData []byte    `db:"encrypted_data"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
