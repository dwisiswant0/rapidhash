package rapidhash

import (
	"math/bits"
	"unsafe"
)

// Secret constants for rapidhash.
// These are carefully chosen constants that provide good hash distribution.
const (
	secret0 = 0x2d358dccaa6c78a5
	secret1 = 0x8bb84b93962eacc9
	secret2 = 0x4b33a62ed433d4a3
	secret3 = 0x4d5a2da51de1aa47
	secret4 = 0xa0761d6478bd642f
	secret5 = 0xe7037ed1a0b428db
	secret6 = 0x90ed1765281c388c
	secret7 = 0xaaaaaaaaaaaaaaaa

	// Precomputed: mix(secret2, secret1) - used when seed=0 to skip a multiply
	seed0Mixed = 0x422765567d8fbfd6
)

// mix performs a 64x64 -> 128 bit multiply, then XORs the high and low 64 bits
// together. This is the core mixing function.
//
//go:inline
func mix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)

	return lo ^ hi
}

// mum performs a 64x64 -> 128 bit multiplication and returns low, high.
//
//go:inline
func mum(a, b uint64) (uint64, uint64) {
	hi, lo := bits.Mul64(a, b)

	return lo, hi
}

// u64 reads a little-endian uint64 using unsafe pointer arithmetic.
// This eliminates bounds checking for maximum performance.
//
//go:nosplit
func u64(p unsafe.Pointer) uint64 {
	return *(*uint64)(p)
}

// u32 reads a little-endian uint32 using unsafe pointer arithmetic.
//
//go:nosplit
func u32(p unsafe.Pointer) uint64 {
	return uint64(*(*uint32)(p))
}

// add returns a pointer offset by n bytes.
//
//go:nosplit
func add(p unsafe.Pointer, n uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + n)
}

// loadUpTo3 reads up to 3 bytes into a single 32-bit lane and extracts the
// first, middle, and last bytes in one unaligned load.
// length must be in [1,3].
func loadUpTo3(data []byte) (uint64, uint64) {
	return (uint64(data[0]) << 45) | uint64(data[len(data)-1]), uint64(data[len(data)>>1])
}

// hashLargeBlock processes blocks >448 bytes using assembly.
// Isolated in its own function to prevent register pressure from affecting
// small inputs.
//
//go:noinline
func hashLargeBlock(p unsafe.Pointer, i int, seed uint64) (unsafe.Pointer, int, uint64) {
	see1, see2 := seed, seed
	see3, see4 := seed, seed
	see5, see6 := seed, seed

	p, i, seed, see1, see2, see3, see4, see5, see6 = accumBlocks(p, i, seed, see1, see2, see3, see4, see5, see6)

	seed ^= see1
	see2 ^= see3
	see4 ^= see5
	seed ^= see6
	see2 ^= see4
	seed ^= see2

	return p, i, seed
}

// Hash computes a 64-bit rapidhash of the input data using the default seed (0).
// This is optimized with a precomputed seed constant to skip one multiply.
func Hash(data []byte) uint64 {
	length := len(data)

	// Empty input fast path
	if length == 0 {
		var a, b uint64 = secret1, seed0Mixed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1)
	}

	p := unsafe.Pointer(unsafe.SliceData(data))

	// Small input fast paths with early returns
	if length <= 16 {
		if length >= 4 {
			var a, b uint64
			if length >= 8 {
				a = u64(p)
				b = u64(add(p, uintptr(length-8)))
			} else {
				a = u32(p)
				b = u32(add(p, uintptr(length-4)))
			}

			// Inline finalization for small inputs
			a ^= secret1
			b ^= seed0Mixed ^ uint64(length)
			a, b = mum(a, b)

			return mix(a^secret7, b^secret1^uint64(length))
		}

		// 1-3 bytes
		a, b := loadUpTo3(data)
		a ^= secret1
		b ^= seed0Mixed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1^uint64(length))
	}

	// For larger inputs, use the full path with precomputed seed
	return hashWithMixedSeed(data, p, length, seed0Mixed)
}

// hashWithMixedSeed is the common path for inputs > 16 bytes.
// seed must already be mixed (seed ^= mix(seed^secret2, secret1)).
//
//go:noinline
func hashWithMixedSeed(data []byte, p unsafe.Pointer, length int, seed uint64) uint64 {
	i := length

	if length > 448 {
		p, i, seed = hashLargeBlock(p, i, seed)
	} else if length > 112 {
		// Inline loop for 113-448 bytes
		see1, see2 := seed, seed
		see3, see4 := seed, seed
		see5, see6 := seed, seed

		for i > 112 {
			seed = mix(u64(p)^secret0, u64(add(p, 8))^seed)
			see1 = mix(u64(add(p, 16))^secret1, u64(add(p, 24))^see1)
			see2 = mix(u64(add(p, 32))^secret2, u64(add(p, 40))^see2)
			see3 = mix(u64(add(p, 48))^secret3, u64(add(p, 56))^see3)
			see4 = mix(u64(add(p, 64))^secret4, u64(add(p, 72))^see4)
			see5 = mix(u64(add(p, 80))^secret5, u64(add(p, 88))^see5)
			see6 = mix(u64(add(p, 96))^secret6, u64(add(p, 104))^see6)
			p = add(p, 112)
			i -= 112
		}

		seed ^= see1
		see2 ^= see3
		see4 ^= see5
		seed ^= see6
		see2 ^= see4
		seed ^= see2
	}

	if i > 16 {
		seed = mix(u64(p)^secret2, u64(add(p, 8))^seed)
		if i > 32 {
			seed = mix(u64(add(p, 16))^secret2, u64(add(p, 24))^seed)
		}
		if i > 48 {
			seed = mix(u64(add(p, 32))^secret1, u64(add(p, 40))^seed)
		}
		if i > 64 {
			seed = mix(u64(add(p, 48))^secret1, u64(add(p, 56))^seed)
		}
		if i > 80 {
			seed = mix(u64(add(p, 64))^secret2, u64(add(p, 72))^seed)
		}
		if i > 96 {
			seed = mix(u64(add(p, 80))^secret1, u64(add(p, 88))^seed)
		}
	}

	// Read from the last 16 and 8 bytes of the ORIGINAL data
	origP := unsafe.Pointer(unsafe.SliceData(data))
	a := u64(add(origP, uintptr(length-16))) ^ uint64(i)
	b := u64(add(origP, uintptr(length-8)))

	a ^= secret1
	b ^= seed
	a, b = mum(a, b)

	return mix(a^secret7, b^secret1^uint64(i))
}

// HashWithSeed computes a 64-bit rapidhash of the input data using the provided
// seed.
func HashWithSeed(data []byte, seed uint64) uint64 {
	length := len(data)
	if length == 0 {
		seed ^= mix(seed^secret2, secret1)

		var a, b uint64 = secret1, seed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1)
	}

	p := unsafe.Pointer(unsafe.SliceData(data))
	seed ^= mix(seed^secret2, secret1)

	// Small input paths with early returns
	if length <= 16 {
		var a, b uint64
		if length >= 4 {
			if length >= 8 {
				a = u64(p)
				b = u64(add(p, uintptr(length-8)))
			} else {
				a = u32(p)
				b = u32(add(p, uintptr(length-4)))
			}

			a ^= secret1
			b ^= seed ^ uint64(length)
			a, b = mum(a, b)

			return mix(a^secret7, b^secret1^uint64(length))
		}

		// 1-3 bytes
		a, b = loadUpTo3(data)
		a ^= secret1
		b ^= seed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1^uint64(length))
	}

	// For larger inputs, use the shared path
	return hashWithMixedSeed(data, p, length, seed)
}

// HashNano computes a hash using the Nano variant, optimized for mobile or
// embedded.
//
// ~13% faster for inputs up to 48 bytes, may be slower for larger inputs.
func HashNano(data []byte) uint64 {
	length := len(data)
	if length == 0 {
		var a, b uint64 = secret1, seed0Mixed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1)
	}

	p := unsafe.Pointer(unsafe.SliceData(data))

	// Small input fast path with early return
	if length <= 16 {
		var a, b uint64
		if length >= 4 {
			if length >= 8 {
				a = u64(p)
				b = u64(add(p, uintptr(length-8)))
			} else {
				a = u32(p)
				b = u32(add(p, uintptr(length-4)))
			}

			a ^= secret1
			b ^= seed0Mixed ^ uint64(length)
			a, b = mum(a, b)

			return mix(a^secret7, b^secret1^uint64(length))
		}

		// 1-3 bytes
		a, b = loadUpTo3(data)
		a ^= secret1
		b ^= seed0Mixed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1^uint64(length))
	}

	return hashNanoWithMixedSeed(data, p, length, seed0Mixed)
}

// HashNanoWithSeed computes a hash using the Nano variant with a custom seed.
func HashNanoWithSeed(data []byte, seed uint64) uint64 {
	length := len(data)
	if length == 0 {
		seed ^= mix(seed^secret2, secret1)

		var a, b uint64 = secret1, seed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1)
	}

	p := unsafe.Pointer(unsafe.SliceData(data))
	seed ^= mix(seed^secret2, secret1)

	// Small input paths with early returns
	if length <= 16 {
		var a, b uint64
		if length >= 4 {
			if length >= 8 {
				a = u64(p)
				b = u64(add(p, uintptr(length-8)))
			} else {
				a = u32(p)
				b = u32(add(p, uintptr(length-4)))
			}

			a ^= secret1
			b ^= seed ^ uint64(length)
			a, b = mum(a, b)

			return mix(a^secret7, b^secret1^uint64(length))
		}

		// 1-3 bytes
		a, b = loadUpTo3(data)
		a ^= secret1
		b ^= seed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1^uint64(length))
	}

	return hashNanoWithMixedSeed(data, p, length, seed)
}

// hashNanoWithMixedSeed handles inputs >16 bytes for HashNano.
//
//go:noinline
func hashNanoWithMixedSeed(data []byte, p unsafe.Pointer, length int, seed uint64) uint64 {
	i := length

	if length > 48 {
		see1, see2 := seed, seed

		for i > 48 {
			seed = mix(u64(p)^secret0, u64(add(p, 8))^seed)
			see1 = mix(u64(add(p, 16))^secret1, u64(add(p, 24))^see1)
			see2 = mix(u64(add(p, 32))^secret2, u64(add(p, 40))^see2)
			p = add(p, 48)
			i -= 48
		}

		seed ^= see1
		seed ^= see2
	}

	if i > 16 {
		seed = mix(u64(p)^secret2, u64(add(p, 8))^seed)
		if i > 32 {
			seed = mix(u64(add(p, 16))^secret2, u64(add(p, 24))^seed)
		}
	}

	origP := unsafe.Pointer(unsafe.SliceData(data))
	a := u64(add(origP, uintptr(length-16))) ^ uint64(i)
	b := u64(add(origP, uintptr(length-8)))

	a ^= secret1
	b ^= seed
	a, b = mum(a, b)

	return mix(a^secret7, b^secret1^uint64(i))
}

// HashMicro computes a hash using the Micro variant, optimized for HPC/server
// applications.
//
// ~16% faster for inputs up to 512 bytes, may be slower for inputs above 1KB.
func HashMicro(data []byte) uint64 {
	length := len(data)
	if length == 0 {
		var a, b uint64 = secret1, seed0Mixed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1)
	}

	p := unsafe.Pointer(unsafe.SliceData(data))

	// Small input fast path with early return
	if length <= 16 {
		var a, b uint64
		if length >= 4 {
			if length >= 8 {
				a = u64(p)
				b = u64(add(p, uintptr(length-8)))
			} else {
				a = u32(p)
				b = u32(add(p, uintptr(length-4)))
			}

			a ^= secret1
			b ^= seed0Mixed ^ uint64(length)
			a, b = mum(a, b)

			return mix(a^secret7, b^secret1^uint64(length))
		}

		// 1-3 bytes
		a, b = loadUpTo3(data)
		a ^= secret1
		b ^= seed0Mixed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1^uint64(length))
	}

	return hashMicroWithMixedSeed(data, p, length, seed0Mixed)
}

// HashMicroWithSeed computes a hash using the Micro variant with a custom seed.
func HashMicroWithSeed(data []byte, seed uint64) uint64 {
	length := len(data)
	if length == 0 {
		seed ^= mix(seed^secret2, secret1)

		var a, b uint64 = secret1, seed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1)
	}

	p := unsafe.Pointer(unsafe.SliceData(data))
	seed ^= mix(seed^secret2, secret1)

	// Small input paths with early returns
	if length <= 16 {
		var a, b uint64
		if length >= 4 {
			if length >= 8 {
				a = u64(p)
				b = u64(add(p, uintptr(length-8)))
			} else {
				a = u32(p)
				b = u32(add(p, uintptr(length-4)))
			}

			a ^= secret1
			b ^= seed ^ uint64(length)
			a, b = mum(a, b)

			return mix(a^secret7, b^secret1^uint64(length))
		}

		// 1-3 bytes
		a, b = loadUpTo3(data)
		a ^= secret1
		b ^= seed
		a, b = mum(a, b)

		return mix(a^secret7, b^secret1^uint64(length))
	}

	return hashMicroWithMixedSeed(data, p, length, seed)
}

// hashMicroWithMixedSeed handles inputs >16 bytes for HashMicro.
//
//go:noinline
func hashMicroWithMixedSeed(data []byte, p unsafe.Pointer, length int, seed uint64) uint64 {
	i := length

	if length > 80 {
		see1, see2 := seed, seed
		see3, see4 := seed, seed

		for i > 80 {
			seed = mix(u64(p)^secret0, u64(add(p, 8))^seed)
			see1 = mix(u64(add(p, 16))^secret1, u64(add(p, 24))^see1)
			see2 = mix(u64(add(p, 32))^secret2, u64(add(p, 40))^see2)
			see3 = mix(u64(add(p, 48))^secret3, u64(add(p, 56))^see3)
			see4 = mix(u64(add(p, 64))^secret4, u64(add(p, 72))^see4)
			p = add(p, 80)
			i -= 80
		}

		seed ^= see1
		see2 ^= see3
		seed ^= see4
		seed ^= see2
	}

	if i > 16 {
		seed = mix(u64(p)^secret2, u64(add(p, 8))^seed)
		if i > 32 {
			seed = mix(u64(add(p, 16))^secret2, u64(add(p, 24))^seed)
		}
		if i > 48 {
			seed = mix(u64(add(p, 32))^secret1, u64(add(p, 40))^seed)
		}
		if i > 64 {
			seed = mix(u64(add(p, 48))^secret1, u64(add(p, 56))^seed)
		}
	}

	origP := unsafe.Pointer(unsafe.SliceData(data))
	a := u64(add(origP, uintptr(length-16))) ^ uint64(i)
	b := u64(add(origP, uintptr(length-8)))

	a ^= secret1
	b ^= seed
	a, b = mum(a, b)

	return mix(a^secret7, b^secret1^uint64(i))
}

// stringToBytes converts a string to a byte slice without allocation.
//
// The returned slice shares memory with the string and must not be modified.
func stringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// HashString computes a 64-bit rapidhash of the input string using the default
// seed (0).
func HashString(s string) uint64 {
	return Hash(stringToBytes(s))
}

// HashStringWithSeed computes a 64-bit rapidhash of the input string using the
// provided seed.
func HashStringWithSeed(s string, seed uint64) uint64 {
	return HashWithSeed(stringToBytes(s), seed)
}

// HashStringNano computes a hash of the input string using the Nano variant.
func HashStringNano(s string) uint64 {
	return HashNano(stringToBytes(s))
}

// HashStringNanoWithSeed computes a hash of the input string using the Nano
// variant with a custom seed.
func HashStringNanoWithSeed(s string, seed uint64) uint64 {
	return HashNanoWithSeed(stringToBytes(s), seed)
}

// HashStringMicro computes a hash of the input string using the Micro variant.
func HashStringMicro(s string) uint64 {
	return HashMicro(stringToBytes(s))
}

// HashStringMicroWithSeed computes a hash of the input string using the Micro
// variant with a custom seed.
func HashStringMicroWithSeed(s string, seed uint64) uint64 {
	return HashMicroWithSeed(stringToBytes(s), seed)
}
