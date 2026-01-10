package rapidhash_test

import (
	"fmt"
	"testing"

	"go.dw1.io/rapidhash"
)

var sizes = []int{8, 16, 32, 64, 128, 256, 512, 1024, 4096, 8192}

// Sink to prevent compiler optimizations
var sink uint64

func makeData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i)
	}
	return data
}

func BenchmarkComputes(b *testing.B) {
	getRapidhashFunction := func(size int) (func([]byte) uint64, string) {
		switch {
		case size <= 48:
			return rapidhash.HashNano, "HashNano"
		case size <= 512:
			return rapidhash.HashMicro, "HashMicro"
		default:
			return rapidhash.Hash, "Hash"
		}
	}

	for _, size := range sizes {
		data := makeData(size)
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			b.Run("Hash", func(b *testing.B) {
				b.SetBytes(int64(size))
				for b.Loop() {
					sink = rapidhash.Hash(data)
				}
			})

			b.Run("Hasher", func(b *testing.B) {
				b.SetBytes(int64(size))

				h := rapidhash.New()
				for b.Loop() {
					h.Reset()
					_, _ = h.Write(data)
					sink = h.Sum64()
				}
			})

			if f, fname := getRapidhashFunction(size); fname != "Hash" {
				b.Run(fname, func(b *testing.B) {
					b.SetBytes(int64(size))

					for b.Loop() {
						sink = f(data)
					}
				})
			}
		})
	}
}
