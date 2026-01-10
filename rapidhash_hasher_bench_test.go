package rapidhash_test

import (
	"testing"

	"go.dw1.io/rapidhash"
)

func BenchmarkHasher1K_Chunked(b *testing.B) {
	data := make([]byte, 1024)
	h := rapidhash.New()
	b.SetBytes(1024)

	for b.Loop() {
		h.Reset()
		_, _ = h.Write(data[:256])
		_, _ = h.Write(data[256:512])
		_, _ = h.Write(data[512:768])
		_, _ = h.Write(data[768:])
		h.Sum64()
	}
}
