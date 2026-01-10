package rapidhash

import "hash"

var _ hash.Hash = (*Hasher)(nil)
var _ hash.Hash64 = (*Hasher)(nil)
var _ hash.Hash32 = (*Hasher)(nil)

// Hasher implements [hash.Hash32] and [hash.Hash64] for streaming hash computation.
//
// Note: For memory-efficiency with large inputs, consider using [Hash] directly.
type Hasher struct {
	seed uint64
	data []byte
}

// New creates a new Hasher with the default seed (0).
func New() *Hasher {
	return NewWithSeed(0)
}

// NewWithSeed creates a new Hasher with the given seed.
func NewWithSeed(seed uint64) *Hasher {
	return &Hasher{
		seed: seed,
		data: make([]byte, 0, 64), // Initial capacity
	}
}

// Reset resets the hasher to its initial state.
func (h *Hasher) Reset() {
	h.data = h.data[:0]
}

// Size returns the number of bytes Sum will return (8 bytes for a 64-bit hash).
func (h *Hasher) Size() int {
	return 8
}

// BlockSize returns the hash's underlying block size.
func (h *Hasher) BlockSize() int {
	return 112
}

// Write adds more data to the running hash.
func (h *Hasher) Write(p []byte) (n int, err error) {
	h.data = append(h.data, p...)

	return len(p), nil
}

// Sum64 returns the current 64-bit hash value.
func (h *Hasher) Sum64() uint64 {
	return HashWithSeed(h.data, h.seed)
}

// Sum32 returns the lower 32 bits of the current hash value.
func (h *Hasher) Sum32() uint32 {
	v := h.Sum64()

	return uint32(v ^ (v >> 32))
}

// Sum appends the current hash to b and returns the resulting slice.
func (h *Hasher) Sum(b []byte) []byte {
	hash := h.Sum64()

	return append(b,
		byte(hash>>56),
		byte(hash>>48),
		byte(hash>>40),
		byte(hash>>32),
		byte(hash>>24),
		byte(hash>>16),
		byte(hash>>8),
		byte(hash),
	)
}
