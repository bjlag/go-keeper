package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"sync"
)

type Cipher struct {
	once        sync.Once
	gcmInstance cipher.AEAD
}

func (c *Cipher) Encrypt(data, key []byte) ([]byte, error) {
	err := c.setup(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, c.gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return c.gcmInstance.Seal(nonce, nonce, data, nil), nil
}

func (c *Cipher) Decrypt(encryptedData, key []byte) ([]byte, error) {
	err := c.setup(key)
	if err != nil {
		return nil, err
	}

	nonceSize := c.gcmInstance.NonceSize()
	nonce, cipheredText := encryptedData[:nonceSize], encryptedData[nonceSize:]

	return c.gcmInstance.Open(nil, nonce, cipheredText, nil)
}

func (c *Cipher) setup(key []byte) error {
	var err error

	c.once.Do(func() {
		var aesBlock cipher.Block
		aesBlock, err = aes.NewCipher(key)
		if err != nil {
			return
		}

		var gcmInstance cipher.AEAD
		gcmInstance, err = cipher.NewGCM(aesBlock)
		if err != nil {
			return
		}

		c.gcmInstance = gcmInstance
	})

	return err
}
