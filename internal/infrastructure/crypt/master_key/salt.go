package master_key

import (
	"crypto/rand"
	"encoding/base64"
)

type Salt []byte

func NewSalt(length int) (*Salt, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	s := Salt(salt)

	return &s, nil
}

func ParseString(salt string) (*Salt, error) {
	b, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	s := Salt(b)

	return &s, nil
}

func (s Salt) Value() []byte {
	return s
}

func (s Salt) ToString() string {
	return base64.StdEncoding.EncodeToString(s)
}
