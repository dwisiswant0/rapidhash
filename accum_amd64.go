//go:build amd64

package rapidhash

import "unsafe"

// accumBlocks processes multiple 112-byte blocks using optimized assembly.
// Returns the new pointer position, remaining length, and updated accumulators.
//
//go:noescape
func accumBlocks(p unsafe.Pointer, length int, seed, see1, see2, see3, see4, see5, see6 uint64) (
	newP unsafe.Pointer, remaining int, nseed, nsee1, nsee2, nsee3, nsee4, nsee5, nsee6 uint64)
