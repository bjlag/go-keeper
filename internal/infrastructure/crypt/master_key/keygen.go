package master_key

import (
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

// KeyGenerator генератор мастер ключа на основе алгоритма PBKDF2.
// В iteCount количество итераций при генерации ключа.
// В keyLen длина ключа.
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

// GenerateMasterKey генерирует мастер ключ для переданного пароля password используя соль salt.
func (g KeyGenerator) GenerateMasterKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, g.iteCount, g.keyLen, sha512.New)
}
