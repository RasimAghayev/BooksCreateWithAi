package shift

import (
	"bytes"
	"crypto/cipher"
)

var cipherCases = []struct {
	key        []byte
	plaintext  []byte
	ciphertext []byte
}{
	{
		key:        bytes.Repeat([]byte{1}, 32),
		plaintext:  bytes.Repeat([]byte("A"), 32),
		ciphertext: bytes.Repeat([]byte("B"), 32),
	},
}

type encrypter struct {
	block     cipher.Block
	blockSize int
}

func NewEncrypter(block cipher.Block) cipher.BlockMode {
	return &encrypter{
		block:     block,
		blockSize: block.BlockSize(),
	}
}

func (e *encrypter) BlockSize() int {
	return e.blockSize
}

func (e *encrypter) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("encrypter: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("encrypter: output smaller than input")
	}
	for len(src) > 0 {
		e.block.Encrypt(dst[:e.blockSize], src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}
