package master_key

import (
	"crypto/rand"
	"encoding/base64"
)

// Salt соль для мастер ключа.
type Salt []byte

// NewSalt создает соль указанной длины length.
func NewSalt(length int) (*Salt, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	s := Salt(salt)

	return &s, nil
}

// ParseString создание структуры Salt из строкового представления соли.
func ParseString(salt string) (*Salt, error) {
	b, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	s := Salt(b)

	return &s, nil
}

// Value возвращает значение соли.
func (s Salt) Value() []byte {
	return s
}

// ToString encoding Salt в строку.
func (s Salt) ToString() string {
	return base64.StdEncoding.EncodeToString(s)
}
