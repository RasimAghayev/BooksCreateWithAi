package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bitfield/shift"
)

func main() {
	keyHex := flag.String("key", "0101010101010101010101010101010101010101010101010101010101010101", "32-byte key in hexadecimal")
	flag.Parse()
	key, err := hexDecode(*keyHex)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	block, err := shift.NewCipher(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dec := shift.NewDecrypter(block)
	plaintext := make([]byte, len(ciphertext))
	dec.CryptBlocks(plaintext, ciphertext)
	os.Stdout.Write(plaintext)
}

func hexDecode(s string) ([]byte, error) {
	data := make([]byte, 0, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		var b byte
		hi := s[i]
		lo := s[i+1]
		switch {
		case hi >= '0' && hi <= '9':
			b = (hi - '0') << 4
		case hi >= 'a' && hi <= 'f':
			b = (hi - 'a' + 10) << 4
		case hi >= 'A' && hi <= 'F':
			b = (hi - 'A' + 10) << 4
		default:
			return nil, fmt.Errorf("invalid hex character: %c", hi)
		}
		switch {
		case lo >= '0' && lo <= '9':
			b |= lo - '0'
		case lo >= 'a' && lo <= 'f':
			b |= lo - 'a' + 10
		case lo >= 'A' && lo <= 'F':
			b |= lo - 'A' + 10
		default:
			return nil, fmt.Errorf("invalid hex character: %c", lo)
		}
		data = append(data, b)
	}
	return data, nil
}
