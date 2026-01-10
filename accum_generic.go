//go:build !amd64

package rapidhash

import "unsafe"

// accumBlocks processes multiple 112-byte blocks.
//
// This is the generic fallback for non-amd64 platforms.
func accumBlocks(p unsafe.Pointer, length int, seed, see1, see2, see3, see4, see5, see6 uint64) (
	newP unsafe.Pointer, remaining int, nseed, nsee1, nsee2, nsee3, nsee4, nsee5, nsee6 uint64) {

	for length > 112 {
		seed = mix(u64(p)^secret0, u64(add(p, 8))^seed)
		see1 = mix(u64(add(p, 16))^secret1, u64(add(p, 24))^see1)
		see2 = mix(u64(add(p, 32))^secret2, u64(add(p, 40))^see2)
		see3 = mix(u64(add(p, 48))^secret3, u64(add(p, 56))^see3)
		see4 = mix(u64(add(p, 64))^secret4, u64(add(p, 72))^see4)
		see5 = mix(u64(add(p, 80))^secret5, u64(add(p, 88))^see5)
		see6 = mix(u64(add(p, 96))^secret6, u64(add(p, 104))^see6)
		p = add(p, 112)
		length -= 112
	}

	return p, length, seed, see1, see2, see3, see4, see5, see6
}
