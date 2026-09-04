package hash

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLenHashReturnsExpectedResult(t *testing.T) {
	t.Parallel()
	input := []byte("I love you, Bob")
	want := []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x0f,
	}
	got := LenHash(input)
	if !bytes.Equal(want, got) {
		t.Errorf("%s: want %x, got %x", input, want, got)
	}
}

func TestSumHashReturnsExpectedResult(t *testing.T) {
	t.Parallel()
	input := []byte("I love you, Bob")
	want := []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x04, 0xfb,
	}
	got := SumHash(input)
	if !bytes.Equal(want, got) {
		t.Errorf("%s: want %x, got %x", input, want, got)
	}
}

func TestLenHashDifferentInputsSameLengthProduceSameHash(t *testing.T) {
	t.Parallel()
	input1 := []byte("I love you, Bob")
	input2 := []byte("I hate you, Bob")
	got1 := LenHash(input1)
	got2 := LenHash(input2)
	if !bytes.Equal(got1, got2) {
		t.Errorf("inputs of same length should produce same hash: got %x and %x", got1, got2)
	}
}

func TestSumHashDifferentInputsProduceDifferentHashes(t *testing.T) {
	t.Parallel()
	input1 := []byte("I love you, Bob")
	input2 := []byte("I hate you, Bob")
	got1 := SumHash(input1)
	got2 := SumHash(input2)
	if bytes.Equal(got1, got2) {
		t.Error("same length inputs should produce different hashes, but got same")
	}
}

func TestSHA256(t *testing.T) {
	t.Parallel()
	input := []byte("hello")
	got := SHA256(input)
	if len(got) != 32 {
		t.Errorf("SHA256 digest should be 32 bytes, got %d", len(got))
	}
}

func TestSumHashAvalancheEffect(t *testing.T) {
	t.Parallel()
	input1 := []byte("I love you, Bob")
	input2 := []byte("I love you, job")
	got1 := SumHash(input1)
	got2 := SumHash(input2)
	if bytes.Equal(got1, got2) {
		t.Error("small input change should cause different hash output")
	}
}

var _ = binary.BigEndian
