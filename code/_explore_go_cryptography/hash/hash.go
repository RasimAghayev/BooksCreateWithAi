package hash

import (
	"crypto/sha256"
	"encoding/binary"
)

func LenHash(input []byte) []byte {
	digest := make([]byte, 8)
	binary.BigEndian.PutUint64(digest, uint64(len(input)))
	return digest
}

func SumHash(input []byte) []byte {
	digest := make([]byte, 8)
	var sum uint64
	for _, b := range input {
		sum += uint64(b)
	}
	binary.BigEndian.PutUint64(digest, sum)
	return digest
}

func SHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
