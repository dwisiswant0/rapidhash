package benchmarks_test

import (
	"fmt"
	"testing"

	"github.com/cespare/xxhash/v2"
	"github.com/dgryski/go-wyhash"
	"go.dw1.io/rapidhash"
)

var sizes = []int{8, 16, 32, 64, 128, 256, 512, 1024, 4096, 8192}

func makeData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i)
	}
	return data
}

// Sink to prevent compiler optimizations
var sink uint64

func BenchmarkXXH64(b *testing.B) {
	for _, size := range sizes {
		data := makeData(size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			for b.Loop() {
				sink = xxhash.Sum64(data)
			}
		})
	}
}

func BenchmarkWyhash(b *testing.B) {
	for _, size := range sizes {
		data := makeData(size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			for b.Loop() {
				sink = wyhash.Hash(data, 0)
			}
		})
	}
}

func BenchmarkRapidhash(b *testing.B) {
	for _, size := range sizes {
		data := makeData(size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			for b.Loop() {
				sink = rapidhash.Hash(data)
			}
		})
	}
}
