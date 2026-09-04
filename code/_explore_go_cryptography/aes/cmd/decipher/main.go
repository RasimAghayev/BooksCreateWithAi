package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bitfield/aes"
)

func main() {
	keyHex := flag.String("key", "", "32-byte key in hexadecimal")
	flag.Parse()
	if *keyHex == "" {
		fmt.Fprintln(os.Stderr, "key flag is required")
		os.Exit(1)
	}
	key := hexDecode(*keyHex)
	if len(key) != 32 {
		fmt.Fprintf(os.Stderr, "key must be 32 bytes (64 hex chars), got %d\n", len(key))
		os.Exit(1)
	}
	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	plaintext, err := aes.Open(key, ciphertext)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Stdout.Write(plaintext)
}

func hexDecode(s string) []byte {
	if len(s)%2 != 0 {
		return nil
	}
	data := make([]byte, len(s)/2)
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
			return nil
		}
		switch {
		case lo >= '0' && lo <= '9':
			b |= lo - '0'
		case lo >= 'a' && lo <= 'f':
			b |= lo - 'a' + 10
		case lo >= 'A' && lo <= 'F':
			b |= lo - 'A' + 10
		default:
			return nil
		}
		data[i/2] = b
	}
	return data
}
