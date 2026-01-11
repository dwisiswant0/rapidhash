package rapidhash_test

import (
	"bytes"
	"hash"
	"testing"

	"go.dw1.io/rapidhash"
)

func TestHasher(t *testing.T) {
	data := []byte("Hello, World!")
	expected := rapidhash.Hash(data)

	h := rapidhash.New()
	_, _ = h.Write(data)
	got := h.Sum64()

	if got != expected {
		t.Errorf("Hasher.Sum64() = 0x%x, want 0x%x", got, expected)
	}
}

func TestHasherChunked(t *testing.T) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	expected := rapidhash.Hash(data)

	// Write in various chunk sizes
	chunkSizes := []int{1, 2, 3, 5, 7, 11, 13, 17}

	for _, chunkSize := range chunkSizes {
		h := rapidhash.New()
		for i := 0; i < len(data); i += chunkSize {
			end := i + chunkSize
			if end > len(data) {
				end = len(data)
			}
			_, _ = h.Write(data[i:end])
		}

		got := h.Sum64()
		if got != expected {
			t.Errorf("Hasher with chunk size %d: Sum64() = 0x%x, want 0x%x",
				chunkSize, got, expected)
		}
	}
}

func TestHasherLargeInput(t *testing.T) {
	// Test with input larger than 112 bytes (block size)
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}
	expected := rapidhash.Hash(data)

	h := rapidhash.New()
	_, _ = h.Write(data)
	got := h.Sum64()

	if got != expected {
		t.Errorf("Hasher with 1000 bytes: Sum64() = 0x%x, want 0x%x", got, expected)
	}
}

func TestHasherReset(t *testing.T) {
	data := []byte("test data")

	h := rapidhash.New()
	_, _ = h.Write(data)
	first := h.Sum64()

	h.Reset()
	_, _ = h.Write(data)
	second := h.Sum64()

	if first != second {
		t.Errorf("After Reset, hashes differ: 0x%x vs 0x%x", first, second)
	}
}

func TestHasherWithSeed(t *testing.T) {
	data := []byte("hello world")

	h0 := rapidhash.New()
	_, _ = h0.Write(data)

	h1 := rapidhash.NewWithSeed(12345)
	_, _ = h1.Write(data)

	if h0.Sum64() == h1.Sum64() {
		t.Error("Different seeds should produce different hashes")
	}

	if h0.Sum64() != rapidhash.Hash(data) {
		t.Error("New() should match Hash()")
	}

	if h1.Sum64() != rapidhash.HashWithSeed(data, 12345) {
		t.Error("NewWithSeed(12345) should match HashWithSeed(data, 12345)")
	}
}

func TestHasherSum(t *testing.T) {
	data := []byte("test")
	h := rapidhash.New()
	_, _ = h.Write(data)

	expected := h.Sum64()
	sum := h.Sum(nil)

	if len(sum) != 8 {
		t.Errorf("Sum() returned %d bytes, want 8", len(sum))
	}

	// Verify the bytes are big-endian encoding of Sum64
	var got uint64
	for i := 0; i < 8; i++ {
		got = (got << 8) | uint64(sum[i])
	}

	if got != expected {
		t.Errorf("Sum() bytes = 0x%x, want 0x%x", got, expected)
	}
}

func TestHasherImplementsHashHash64(t *testing.T) {
	var _ hash.Hash64 = rapidhash.New()
}

func TestHasherSize(t *testing.T) {
	h := rapidhash.New()
	if h.Size() != 8 {
		t.Errorf("Size() = %d, want 8", h.Size())
	}
}

func TestHasherBlockSize(t *testing.T) {
	h := rapidhash.New()
	if h.BlockSize() != 112 {
		t.Errorf("BlockSize() = %d, want 112", h.BlockSize())
	}
}

func TestHasherConsistency(t *testing.T) {
	data := bytes.Repeat([]byte("abcdefghij"), 100)

	var results [10]uint64
	for i := range results {
		h := rapidhash.New()
		_, _ = h.Write(data)
		results[i] = h.Sum64()
	}

	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Inconsistent results: results[%d] = 0x%x, results[0] = 0x%x",
				i, results[i], results[0])
		}
	}
}

func TestHasherMatchesHash(t *testing.T) {
	sizes := []int{0, 1, 7, 8, 15, 16, 17, 32, 48, 64, 80, 100, 112, 113, 150, 200, 500, 1000}

	for _, size := range sizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		expected := rapidhash.Hash(data)

		// Test writing all at once
		h := rapidhash.New()
		_, _ = h.Write(data)
		if got := h.Sum64(); got != expected {
			t.Errorf("size=%d (all at once): Hasher.Sum64() = 0x%x, want 0x%x", size, got, expected)
		}

		// Test writing in 1-byte chunks
		h = rapidhash.New()
		for i := 0; i < size; i++ {
			_, _ = h.Write(data[i : i+1])
		}
		if got := h.Sum64(); got != expected {
			t.Errorf("size=%d (1-byte chunks): Hasher.Sum64() = 0x%x, want 0x%x", size, got, expected)
		}

		// Test writing in 7-byte chunks (non-aligned)
		h = rapidhash.New()
		for i := 0; i < size; i += 7 {
			end := i + 7
			if end > size {
				end = size
			}
			_, _ = h.Write(data[i:end])
		}
		if got := h.Sum64(); got != expected {
			t.Errorf("size=%d (7-byte chunks): Hasher.Sum64() = 0x%x, want 0x%x", size, got, expected)
		}
	}
}

func TestHasherWriteString(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"hello",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 200)),
	}

	for _, s := range testStrings {
		// WriteString should produce same result as Write([]byte(s))
		h1 := rapidhash.New()
		_, _ = h1.Write([]byte(s))

		h2 := rapidhash.New()
		_, _ = h2.WriteString(s)

		if h1.Sum64() != h2.Sum64() {
			t.Errorf("WriteString(%q) = 0x%x, Write([]byte) = 0x%x", s, h2.Sum64(), h1.Sum64())
		}

		// Should also match HashString
		if h2.Sum64() != rapidhash.HashString(s) {
			t.Errorf("WriteString(%q) = 0x%x, HashString = 0x%x", s, h2.Sum64(), rapidhash.HashString(s))
		}
	}
}

func TestHasherWriteStringChunked(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog"
	expected := rapidhash.HashString(s)

	// Write in multiple string chunks
	h := rapidhash.New()
	_, _ = h.WriteString("The quick ")
	_, _ = h.WriteString("brown fox ")
	_, _ = h.WriteString("jumps over ")
	_, _ = h.WriteString("the lazy dog")

	if got := h.Sum64(); got != expected {
		t.Errorf("WriteString chunked = 0x%x, want 0x%x", got, expected)
	}
}

func TestHasherWriteStringMixed(t *testing.T) {
	// Mix Write and WriteString calls
	h1 := rapidhash.New()
	_, _ = h1.Write([]byte("hello "))
	_, _ = h1.WriteString("world")

	h2 := rapidhash.New()
	_, _ = h2.WriteString("hello ")
	_, _ = h2.Write([]byte("world"))

	h3 := rapidhash.New()
	_, _ = h3.Write([]byte("hello world"))

	if h1.Sum64() != h2.Sum64() || h2.Sum64() != h3.Sum64() {
		t.Errorf("Mixed Write/WriteString produced different results: 0x%x, 0x%x, 0x%x",
			h1.Sum64(), h2.Sum64(), h3.Sum64())
	}
}
