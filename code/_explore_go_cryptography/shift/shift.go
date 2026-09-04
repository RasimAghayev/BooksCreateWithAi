package shift

import (
	"bytes"
	"errors"
)

const BlockSize = 32

var testKey = bytes.Repeat([]byte{1}, BlockSize)

func Encipher(plaintext []byte, key []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	for i, b := range plaintext {
		ciphertext[i] = b + key[i%len(key)]
	}
	return ciphertext
}

func Decipher(ciphertext []byte, key []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	for i, b := range ciphertext {
		plaintext[i] = b - key[i%len(key)]
	}
	return plaintext
}

func Crack(ciphertext, crib []byte) ([]byte, error) {
	const MaxKeyLen = 32
	cribLen := min(len(crib), len(ciphertext))
	maxKeyLen := min(MaxKeyLen, cribLen/2)

	for keyLen := 1; keyLen <= maxKeyLen; keyLen++ {
		key, ok := crackKey(ciphertext[:cribLen], crib, keyLen)
		if !ok {
			continue
		}
		verified := true
		for i := keyLen; i < cribLen; i++ {
			if ciphertext[i]-key[i%keyLen] != crib[i] {
				verified = false
				break
			}
		}
		if !verified {
			continue
		}
		if bytes.Equal(crib, Decipher(ciphertext[:cribLen], key)) {
			return key, nil
		}
	}
	return nil, errors.New("no key found")
}

func crackKey(ciphertext, crib []byte, keyLen int) ([]byte, bool) {
	key := make([]byte, keyLen)
	for k := range keyLen {
		found := false
		for guess := range 256 {
			guessByte := byte(guess)
			match := true
			for i := k; i < len(crib); i += keyLen {
				if ciphertext[i]-guessByte != crib[i] {
					match = false
					break
				}
			}
			if match {
				if !found {
					key[k] = guessByte
					found = true
				} else if key[k] != guessByte {
					return nil, false
				}
			}
		}
		if !found {
			return nil, false
		}
	}
	return key, true
}
