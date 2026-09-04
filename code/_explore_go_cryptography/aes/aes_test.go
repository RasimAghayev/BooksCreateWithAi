package aes

import (
	"bytes"
	cryptoaes "crypto/aes"
	"testing"
)

func TestPad(t *testing.T) {
	t.Parallel()
	input := []byte("SHORT")
	want := []byte("SHORT\013\013\013\013\013\013\013\013\013\013\013")
	got := Pad(input, 16)
	if !bytes.Equal(want, got) {
		t.Errorf("want %x, got %x", want, got)
	}
}

func TestUnpad(t *testing.T) {
	t.Parallel()
	padded := []byte("SHORT\013\013\013\013\013\013\013\013\013\013\013")
	want := []byte("SHORT")
	got, err := Unpad(padded, 16)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestEncryptDecryptCBC(t *testing.T) {
	t.Parallel()
	key := bytes.Repeat([]byte{1}, 32)
	plaintext := []byte("The tiger appears at its own pleasure.")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(plaintext, ciphertext[cryptoaes.BlockSize:]) {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("round trip failed: want %q, got %q", plaintext, decrypted)
	}
}

func TestSealOpenGCM(t *testing.T) {
	t.Parallel()
	key := bytes.Repeat([]byte{1}, 32)
	plaintext := []byte("Hello, world!")

	ciphertext, err := Seal(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := Open(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("round trip failed: want %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptWithWrongKeyFails(t *testing.T) {
	t.Parallel()
	key1 := bytes.Repeat([]byte{1}, 32)
	key2 := bytes.Repeat([]byte{2}, 32)
	plaintext := []byte("secret message")

	ciphertext, err := Encrypt(key1, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decrypt(key2, ciphertext)
	if err == nil {
		t.Fatal("expected error decrypting with wrong key, got nil")
	}
}
