package shift

import (
	"crypto/cipher"
	"errors"
)

type shiftCipher struct {
	key [BlockSize]byte
}

func NewCipher(key []byte) (cipher.Block, error) {
	if len(key) != BlockSize {
		return nil, errors.New("incorrect key size")
	}
	return &shiftCipher{
		key: [BlockSize]byte(key),
	}, nil
}

func (c *shiftCipher) BlockSize() int {
	return BlockSize
}

func (c *shiftCipher) Encrypt(dst, src []byte) {
	for i, b := range src {
		dst[i] = b + c.key[i%len(c.key)]
	}
}

func (c *shiftCipher) Decrypt(dst, src []byte) {
	for i, b := range src {
		dst[i] = b - c.key[i%len(c.key)]
	}
}
