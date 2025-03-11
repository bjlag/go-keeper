package master_key

import (
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

type KeyGenerator struct {
	iteCount int
	keyLen   int
}

func NewKeyGenerator(iterCount, keyLen int) *KeyGenerator {
	return &KeyGenerator{
		iteCount: iterCount,
		keyLen:   keyLen,
	}
}

func (g KeyGenerator) GenerateMasterKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, g.iteCount, g.keyLen, sha512.New)
}
