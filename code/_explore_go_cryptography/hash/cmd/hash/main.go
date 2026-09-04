package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bitfield/hash"
)

func main() {
	algo := flag.String("algo", "len", "hash algorithm (len, sum, sha256)")
	flag.Parse()
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var digest []byte
	switch *algo {
	case "len":
		digest = hash.LenHash(data)
	case "sum":
		digest = hash.SumHash(data)
	case "sha256":
		digest = hash.SHA256(data)
	default:
		fmt.Fprintf(os.Stderr, "unknown algorithm: %s\n", *algo)
		os.Exit(1)
	}
	fmt.Printf("%x\n", digest)
}
