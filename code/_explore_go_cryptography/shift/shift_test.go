package shift

import (
	"bytes"
	"crypto/cipher"
	"testing"
)

var cases = []struct {
	key        []byte
	plaintext  []byte
	ciphertext []byte
}{
	{
		key:        []byte{1},
		plaintext:  []byte("HAL"),
		ciphertext: []byte("IBM"),
	},
	{
		key:        []byte{2},
		plaintext:  []byte("SPEC"),
		ciphertext: []byte("URGE"),
	},
	{
		key:        []byte{3},
		plaintext:  []byte("PERK"),
		ciphertext: []byte("SHUN"),
	},
	{
		key:        []byte{4},
		plaintext:  []byte("GEL"),
		ciphertext: []byte("KIP"),
	},
	{
		key:        []byte{7},
		plaintext:  []byte("CHEER"),
		ciphertext: []byte("JOLLY"),
	},
	{
		key:        []byte{10},
		plaintext:  []byte("BEEF"),
		ciphertext: []byte("LOOP"),
	},
	{
		key:        []byte{1, 2},
		plaintext:  []byte{0, 1, 2, 3, 255},
		ciphertext: []byte{1, 3, 3, 5, 0},
	},
}

func TestEncipher(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		name := string(tc.plaintext)
		t.Run(name, func(t *testing.T) {
			got := Encipher(tc.plaintext, tc.key)
			if !bytes.Equal(tc.ciphertext, got) {
				t.Errorf("want %q, got %q", tc.ciphertext, got)
			}
		})
	}
}

func TestDecipher(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		name := string(tc.ciphertext)
		t.Run(name, func(t *testing.T) {
			got := Decipher(tc.ciphertext, tc.key)
			if !bytes.Equal(tc.plaintext, got) {
				t.Errorf("want %q, got %q", tc.plaintext, got)
			}
		})
	}
}

func TestCrack(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		t.Run(string(tc.plaintext), func(t *testing.T) {
			crib := tc.plaintext[:3]
			if len(tc.key) > 1 {
				crib = tc.plaintext
			}
			got, err := Crack(tc.ciphertext, crib)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(tc.key, got) {
				t.Fatalf("want %x, got %x", tc.key, got)
			}
		})
	}
}

func TestCrackReturnsErrorWhenKeyNotFound(t *testing.T) {
	t.Parallel()
	_, err := Crack([]byte("abcdefgh"), []byte("zzzz"))
	if err == nil {
		t.Fatal("want error when key not found, got nil")
	}
}

func TestNewCipherReturnsErrorForWrongKeySize(t *testing.T) {
	t.Parallel()
	_, err := NewCipher([]byte{1, 2, 3})
	if err == nil {
		t.Fatal("want error for wrong key size, got nil")
	}
}

func TestBlockSize(t *testing.T) {
	t.Parallel()
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	want := BlockSize
	got := block.BlockSize()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestEncrypt(t *testing.T) {
	t.Parallel()
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range cipherCases {
		t.Run(string(tc.plaintext), func(t *testing.T) {
			got := make([]byte, len(tc.plaintext))
			block.Encrypt(got, tc.plaintext)
			if !bytes.Equal(tc.ciphertext, got) {
				t.Errorf("want %x, got %x", tc.ciphertext, got)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range cipherCases {
		t.Run(string(tc.ciphertext), func(t *testing.T) {
			got := make([]byte, len(tc.ciphertext))
			block.Decrypt(got, tc.ciphertext)
			if !bytes.Equal(tc.plaintext, got) {
				t.Errorf("want %x, got %x", tc.plaintext, got)
			}
		})
	}
}

func TestEncrypterCorrectlyReportsCipherBlockSize(t *testing.T) {
	t.Parallel()
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	enc := NewEncrypter(block)
	want := BlockSize
	got := enc.BlockSize()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestEncrypterEnciphersBlockAlignedMessage(t *testing.T) {
	t.Parallel()
	plaintext := []byte("This message is exactly 32 bytes")
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	enc := NewEncrypter(block)
	want := []byte("Uijt!nfttbhf!jt!fybdumz!43!czuft")
	got := make([]byte, 32)
	enc.CryptBlocks(got, plaintext)
	if !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestDecrypterDeciphersBlockAlignedMessage(t *testing.T) {
	t.Parallel()
	ciphertext := []byte("Uijt!nfttbhf!jt!fybdumz!43!czuft")
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	dec := NewDecrypter(block)
	want := []byte("This message is exactly 32 bytes")
	got := make([]byte, 32)
	dec.CryptBlocks(got, ciphertext)
	if !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestEncrypterPanicsOnUnalignedInput(t *testing.T) {
	t.Parallel()
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	enc := NewEncrypter(block)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unaligned input, but none occurred")
		}
	}()
	enc.CryptBlocks(make([]byte, 31), make([]byte, 31))
}

func TestCryptBlocksRoundTrip(t *testing.T) {
	t.Parallel()
	plaintext := []byte("This message is exactly 32 bytes")
	block, err := NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	enc := NewEncrypter(block)
	ciphertext := make([]byte, 32)
	enc.CryptBlocks(ciphertext, plaintext)

	dec := NewDecrypter(block)
	got := make([]byte, 32)
	dec.CryptBlocks(got, ciphertext)

	if !bytes.Equal(plaintext, got) {
		t.Errorf("round trip failed: want %q, got %q", plaintext, got)
	}
}

var _ cipher.Block = (*shiftCipher)(nil)
