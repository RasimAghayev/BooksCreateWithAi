package shift

import (
	"crypto/cipher"
)

type decrypter struct {
	block     cipher.Block
	blockSize int
}

func NewDecrypter(block cipher.Block) cipher.BlockMode {
	return &decrypter{
		block:     block,
		blockSize: block.BlockSize(),
	}
}

func (d *decrypter) BlockSize() int {
	return d.blockSize
}

func (d *decrypter) CryptBlocks(dst, src []byte) {
	if len(src)%d.blockSize != 0 {
		panic("decrypter: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("decrypter: output smaller than input")
	}
	for len(src) > 0 {
		d.block.Decrypt(dst[:d.blockSize], src[:d.blockSize])
		src = src[d.blockSize:]
		dst = dst[d.blockSize:]
	}
}
