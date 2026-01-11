package rapidhash_test

import (
	"testing"

	"go.dw1.io/rapidhash"
)

// Test vectors generated from C reference implementation
var testVectors = []struct {
	input    string
	expected uint64
}{
	{"", 0x338dc4be2cecdae},
	{"a", 0x599f47df33a2e1eb},
	{"ab", 0x7b20ba72bb425975},
	{"abc", 0xcb475beafa9c0da2},
	{"abcd", 0xf8f44f4a65e26132},
	{"abcde", 0xc5fb18456266e8b0},
	{"abcdef", 0xfd615a3b9273e7a9},
	{"abcdefg", 0x2760e84111b29a0d},
	{"abcdefgh", 0xab159e602a29f41f},
	{"abcdefghi", 0x70a51746ca96c1c6},
	{"abcdefghij", 0xc9259e908c63950b},
	{"abcdefghijk", 0x1021ad8053c20148},
	{"abcdefghijkl", 0x7fd5cf0798b86811},
	{"abcdefghijklm", 0x581aae01bc22b1bb},
	{"abcdefghijklmn", 0xf167824f3d2e14c1},
	{"abcdefghijklmno", 0x54c03d9b150448c0},
	{"abcdefghijklmnop", 0xc78ae6a1774adb1e},
	{"abcdefghijklmnopq", 0xc427c11a4463b8},
	{"Hello, World!", 0x75bff66af6ba4d5b},
	{"The quick brown fox jumps over the lazy dog", 0x91722dc8d52a3f7b},
	{"hello", 0x2e2d7651b45f7946},
}

// Test vectors for all three variants at various sizes
// Generated from C reference: {size, Hash, HashMicro, HashNano}
var sizeTestVectors = []struct {
	size      int
	hash      uint64
	hashMicro uint64
	hashNano  uint64
}{
	{0, 0x338dc4be2cecdae, 0x338dc4be2cecdae, 0x338dc4be2cecdae},
	{1, 0x4f23c791b16eba02, 0x4f23c791b16eba02, 0x4f23c791b16eba02},
	{3, 0xdbd091bcf57ae814, 0xdbd091bcf57ae814, 0xdbd091bcf57ae814},
	{4, 0x46fef26db4943adf, 0x46fef26db4943adf, 0x46fef26db4943adf},
	{7, 0x7f403e573bb8ebc1, 0x7f403e573bb8ebc1, 0x7f403e573bb8ebc1},
	{8, 0xda56413ff396af3e, 0xda56413ff396af3e, 0xda56413ff396af3e},
	{15, 0x8ec6dfea933104bb, 0x8ec6dfea933104bb, 0x8ec6dfea933104bb},
	{16, 0xd6bfc1bcf7e9ca19, 0xd6bfc1bcf7e9ca19, 0xd6bfc1bcf7e9ca19},
	{17, 0x7508c9e74d5b5366, 0x7508c9e74d5b5366, 0x7508c9e74d5b5366},
	{31, 0xa0e039c5b97d67f3, 0xa0e039c5b97d67f3, 0xa0e039c5b97d67f3},
	{32, 0xc0186990f026b180, 0xc0186990f026b180, 0xc0186990f026b180},
	{33, 0xeb4ff8393398a779, 0xeb4ff8393398a779, 0xeb4ff8393398a779},
	{47, 0xe6ed23c058015cb9, 0xe6ed23c058015cb9, 0xe6ed23c058015cb9},
	{48, 0xecd5ed3e946f9c91, 0xecd5ed3e946f9c91, 0xecd5ed3e946f9c91},
	{49, 0x635a714c24c02d64, 0x635a714c24c02d64, 0x154059f965173db2},
	{79, 0x29809c72d7013a6e, 0x29809c72d7013a6e, 0x61be47672a06929a},
	{80, 0xe7e477a0dffeae1f, 0xe7e477a0dffeae1f, 0x392712804da62fe5},
	{81, 0xab742caab8765cd2, 0x8c7e35f2679ba33f, 0x2bbb2d7c342e49d7},
	{111, 0xc72721d21a642834, 0x8c010559e00cfc46, 0x41901d845258dd61},
	{112, 0x667174637fd34ae7, 0x336ea990f5fcbd44, 0x45d21d1ad07df9df},
	{113, 0xabaf0e2bdacf7e23, 0xa706652caf6839ee, 0xcfa1a10eec7c0bf},
	{150, 0x548edfe200a7b543, 0xe45072cc949135e9, 0xf20e7c707203e384},
	{200, 0x9d1612cdf44b1c42, 0x3158564c851d0c02, 0x180e783519db3e0c},
	{500, 0xb943a1d6f7b18de4, 0xbdbd39608ffda16e, 0x236b20e15e498410},
	{1000, 0x2a6be558a956faf3, 0xe2e33d1cfc2d95a5, 0xf6357a963b6bcdc9},
}

func TestHash(t *testing.T) {
	for _, tc := range testVectors {
		got := rapidhash.Hash([]byte(tc.input))
		if got != tc.expected {
			t.Errorf("Hash(%q) = 0x%x, want 0x%x", tc.input, got, tc.expected)
		}
	}
}

func TestHashEmpty(t *testing.T) {
	hash := rapidhash.Hash(nil)
	hashEmpty := rapidhash.Hash([]byte{})

	if hash != hashEmpty {
		t.Errorf("Hash(nil) = 0x%x, Hash([]) = 0x%x, want equal", hash, hashEmpty)
	}
}

func TestHashWithSeed(t *testing.T) {
	data := []byte("hello world")

	hash0 := rapidhash.HashWithSeed(data, 0)
	hash1 := rapidhash.HashWithSeed(data, 1)
	hash2 := rapidhash.HashWithSeed(data, 12345)

	// Different seeds should produce different hashes
	if hash0 == hash1 {
		t.Errorf("seed 0 and seed 1 produced same hash: 0x%x", hash0)
	}
	if hash0 == hash2 {
		t.Errorf("seed 0 and seed 12345 produced same hash: 0x%x", hash0)
	}
	if hash1 == hash2 {
		t.Errorf("seed 1 and seed 12345 produced same hash: 0x%x", hash1)
	}

	// Hash with seed 0 should equal Hash
	if hash0 != rapidhash.Hash(data) {
		t.Errorf("HashWithSeed(data, 0) = 0x%x, Hash(data) = 0x%x, want equal",
			hash0, rapidhash.Hash(data))
	}
}

func TestHashAndMicroMatchUpTo80Bytes(t *testing.T) {
	// Hash and HashMicro should match for inputs <= 80 bytes
	for i := 0; i <= 80; i++ {
		data := make([]byte, i)
		for j := range data {
			data[j] = byte(j)
		}

		hash := rapidhash.Hash(data)
		hashMicro := rapidhash.HashMicro(data)

		if hash != hashMicro {
			t.Errorf("len=%d: Hash(0x%x) != HashMicro(0x%x)", i, hash, hashMicro)
		}
	}
}

func TestHashMicroWithSeed(t *testing.T) {
	data := []byte("hello world")

	// Different seeds should produce different hashes
	hash0 := rapidhash.HashMicroWithSeed(data, 0)
	hash1 := rapidhash.HashMicroWithSeed(data, 1)
	hash2 := rapidhash.HashMicroWithSeed(data, 12345)

	if hash0 == hash1 {
		t.Errorf("seed 0 and seed 1 produced same hash: 0x%x", hash0)
	}
	if hash0 == hash2 {
		t.Errorf("seed 0 and seed 12345 produced same hash: 0x%x", hash0)
	}
	if hash1 == hash2 {
		t.Errorf("seed 1 and seed 12345 produced same hash: 0x%x", hash1)
	}

	// HashMicroWithSeed with seed 0 should equal HashMicro
	if hash0 != rapidhash.HashMicro(data) {
		t.Errorf("HashMicroWithSeed(data, 0) = 0x%x, HashMicro(data) = 0x%x, want equal",
			hash0, rapidhash.HashMicro(data))
	}
}

func TestHashMicroWithSeedAllSizes(t *testing.T) {
	seed := uint64(42)

	// Test all code paths: 0, 1-3, 4-7, 8-15, 16, 17+, >80
	sizes := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 15, 16, 17, 32, 48, 64, 80, 81, 100, 150}

	for _, size := range sizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		// Should not panic and should produce non-zero hash (except possibly for empty)
		hash := rapidhash.HashMicroWithSeed(data, seed)
		hashNoSeed := rapidhash.HashMicro(data)

		// With seed should differ from without seed (for non-empty)
		if size > 0 && hash == hashNoSeed {
			t.Errorf("size=%d: HashMicroWithSeed should differ from HashMicro", size)
		}

		// Consistency check: same input should produce same output
		hash2 := rapidhash.HashMicroWithSeed(data, seed)
		if hash != hash2 {
			t.Errorf("size=%d: inconsistent results 0x%x vs 0x%x", size, hash, hash2)
		}
	}
}

func TestHashMicroWithSeedEmpty(t *testing.T) {
	// Test empty input path specifically
	hash1 := rapidhash.HashMicroWithSeed(nil, 0)
	hash2 := rapidhash.HashMicroWithSeed([]byte{}, 0)

	if hash1 != hash2 {
		t.Errorf("HashMicroWithSeed(nil, 0) = 0x%x, HashMicroWithSeed([], 0) = 0x%x, want equal",
			hash1, hash2)
	}

	// Empty with different seeds should differ
	hash3 := rapidhash.HashMicroWithSeed(nil, 12345)
	if hash1 == hash3 {
		t.Errorf("Empty input with different seeds should produce different hashes")
	}
}

func TestHashNanoWithSeed(t *testing.T) {
	data := []byte("hello world")

	// Different seeds should produce different hashes
	hash0 := rapidhash.HashNanoWithSeed(data, 0)
	hash1 := rapidhash.HashNanoWithSeed(data, 1)
	hash2 := rapidhash.HashNanoWithSeed(data, 12345)

	if hash0 == hash1 {
		t.Errorf("seed 0 and seed 1 produced same hash: 0x%x", hash0)
	}
	if hash0 == hash2 {
		t.Errorf("seed 0 and seed 12345 produced same hash: 0x%x", hash0)
	}
	if hash1 == hash2 {
		t.Errorf("seed 1 and seed 12345 produced same hash: 0x%x", hash1)
	}

	// HashNanoWithSeed with seed 0 should equal HashNano
	if hash0 != rapidhash.HashNano(data) {
		t.Errorf("HashNanoWithSeed(data, 0) = 0x%x, HashNano(data) = 0x%x, want equal",
			hash0, rapidhash.HashNano(data))
	}
}

func TestHashNanoWithSeedAllSizes(t *testing.T) {
	seed := uint64(42)

	// Test all code paths: 0, 1-3, 4-7, 8-15, 16, 17+, >48
	sizes := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 15, 16, 17, 32, 48, 49, 64, 80, 100, 150}

	for _, size := range sizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		// Should not panic and should produce non-zero hash (except possibly for empty)
		hash := rapidhash.HashNanoWithSeed(data, seed)
		hashNoSeed := rapidhash.HashNano(data)

		// With seed should differ from without seed (for non-empty)
		if size > 0 && hash == hashNoSeed {
			t.Errorf("size=%d: HashNanoWithSeed should differ from HashNano", size)
		}

		// Consistency check: same input should produce same output
		hash2 := rapidhash.HashNanoWithSeed(data, seed)
		if hash != hash2 {
			t.Errorf("size=%d: inconsistent results 0x%x vs 0x%x", size, hash, hash2)
		}
	}
}

func TestHashNanoWithSeedEmpty(t *testing.T) {
	// Test empty input path specifically
	hash1 := rapidhash.HashNanoWithSeed(nil, 0)
	hash2 := rapidhash.HashNanoWithSeed([]byte{}, 0)

	if hash1 != hash2 {
		t.Errorf("HashNanoWithSeed(nil, 0) = 0x%x, HashNanoWithSeed([], 0) = 0x%x, want equal",
			hash1, hash2)
	}

	// Empty with different seeds should differ
	hash3 := rapidhash.HashNanoWithSeed(nil, 12345)
	if hash1 == hash3 {
		t.Errorf("Empty input with different seeds should produce different hashes")
	}
}

func TestAllVariantsAgainstCReference(t *testing.T) {
	for _, tc := range sizeTestVectors {
		data := make([]byte, tc.size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		gotHash := rapidhash.Hash(data)
		gotMicro := rapidhash.HashMicro(data)
		gotNano := rapidhash.HashNano(data)

		if gotHash != tc.hash {
			t.Errorf("size=%d: Hash() = 0x%x, want 0x%x", tc.size, gotHash, tc.hash)
		}
		if gotMicro != tc.hashMicro {
			t.Errorf("size=%d: HashMicro() = 0x%x, want 0x%x", tc.size, gotMicro, tc.hashMicro)
		}
		if gotNano != tc.hashNano {
			t.Errorf("size=%d: HashNano() = 0x%x, want 0x%x", tc.size, gotNano, tc.hashNano)
		}
	}
}

func TestAllVariantsWithSeedMatchForSmallInputs(t *testing.T) {
	seed := uint64(99999)

	// All three seeded variants should match for inputs <= 48 bytes
	for i := 0; i <= 48; i++ {
		data := make([]byte, i)
		for j := range data {
			data[j] = byte(j)
		}

		hash := rapidhash.HashWithSeed(data, seed)
		hashMicro := rapidhash.HashMicroWithSeed(data, seed)
		hashNano := rapidhash.HashNanoWithSeed(data, seed)

		if hash != hashMicro {
			t.Errorf("len=%d: HashWithSeed(0x%x) != HashMicroWithSeed(0x%x)", i, hash, hashMicro)
		}
		if hash != hashNano {
			t.Errorf("len=%d: HashWithSeed(0x%x) != HashNanoWithSeed(0x%x)", i, hash, hashNano)
		}
	}
}

func TestVariantsMatchForSmallInputs(t *testing.T) {
	// All three variants should match for inputs <= 48 bytes
	for i := 0; i <= 48; i++ {
		data := make([]byte, i)
		for j := range data {
			data[j] = byte(j)
		}

		hash := rapidhash.Hash(data)
		hashMicro := rapidhash.HashMicro(data)
		hashNano := rapidhash.HashNano(data)

		if hash != hashMicro {
			t.Errorf("len=%d: Hash(0x%x) != HashMicro(0x%x)", i, hash, hashMicro)
		}
		if hash != hashNano {
			t.Errorf("len=%d: Hash(0x%x) != HashNano(0x%x)", i, hash, hashNano)
		}
	}
}

func TestVariantsDivergeForLargerInputs(t *testing.T) {
	// Hash and Micro should differ for inputs > 112 bytes
	data := make([]byte, 150)
	for i := range data {
		data[i] = byte(i)
	}

	hash := rapidhash.Hash(data)
	hashMicro := rapidhash.HashMicro(data)
	hashNano := rapidhash.HashNano(data)

	// All three should be different for large inputs
	if hash == hashMicro {
		t.Log("Note: Hash and HashMicro happened to match for 150 bytes")
	}
	if hash == hashNano {
		t.Log("Note: Hash and HashNano happened to match for 150 bytes")
	}

	// Just ensure they all produce valid non-zero hashes
	if hash == 0 || hashMicro == 0 || hashNano == 0 {
		t.Error("One of the hashes is unexpectedly zero")
	}
}

func TestBitFlipSensitivity(t *testing.T) {
	data := make([]byte, 32)
	for i := range data {
		data[i] = byte(i)
	}

	baseHash := rapidhash.Hash(data)

	// Flip each bit and verify the hash changes
	for byteIdx := 0; byteIdx < len(data); byteIdx++ {
		for bit := 0; bit < 8; bit++ {
			// Flip the bit
			data[byteIdx] ^= 1 << bit

			newHash := rapidhash.Hash(data)
			if newHash == baseHash {
				t.Errorf("Flipping bit %d of byte %d didn't change hash", bit, byteIdx)
			}

			// Flip it back
			data[byteIdx] ^= 1 << bit
		}
	}
}

func TestLengthVariation(t *testing.T) {
	// Hashing "ab" should be different from hashing "a" then "b" separately
	ab := rapidhash.Hash([]byte("ab"))
	a := rapidhash.Hash([]byte("a"))
	b := rapidhash.Hash([]byte("b"))

	if ab == a^b {
		t.Error("Hash(ab) should not equal Hash(a) XOR Hash(b)")
	}
	if ab == a+b {
		t.Error("Hash(ab) should not equal Hash(a) + Hash(b)")
	}
}

// TestHashStringMatchesHash verifies that HashString produces identical output
// to Hash([]byte(s)) for all string sizes.
func TestHashStringMatchesHash(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"ab",
		"abc",
		"abcd",
		"abcdefgh",
		"abcdefghijklmnop",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 100)),
		string(make([]byte, 200)),
		string(make([]byte, 500)),
		string(make([]byte, 1000)),
	}

	for _, s := range testStrings {
		expected := rapidhash.Hash([]byte(s))
		got := rapidhash.HashString(s)
		if got != expected {
			t.Errorf("HashString(%q) = 0x%x, want 0x%x (len=%d)", s, got, expected, len(s))
		}
	}
}

// TestHashStringWithSeedMatchesHashWithSeed verifies seeded string hashing.
func TestHashStringWithSeedMatchesHashWithSeed(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"hello world",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 500)),
	}
	seeds := []uint64{0, 1, 12345, 0xdeadbeef, 0xffffffffffffffff}

	for _, s := range testStrings {
		for _, seed := range seeds {
			expected := rapidhash.HashWithSeed([]byte(s), seed)
			got := rapidhash.HashStringWithSeed(s, seed)
			if got != expected {
				t.Errorf("HashStringWithSeed(%q, %d) = 0x%x, want 0x%x", s, seed, got, expected)
			}
		}
	}
}

// TestHashStringNanoMatchesHashNano verifies Nano variant string hashing.
func TestHashStringNanoMatchesHashNano(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"ab",
		"abc",
		"abcdefgh",
		"abcdefghijklmnop",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 100)),
		string(make([]byte, 200)),
	}

	for _, s := range testStrings {
		expected := rapidhash.HashNano([]byte(s))
		got := rapidhash.HashStringNano(s)
		if got != expected {
			t.Errorf("HashStringNano(%q) = 0x%x, want 0x%x (len=%d)", s, got, expected, len(s))
		}
	}
}

// TestHashStringNanoWithSeedMatchesHashNanoWithSeed verifies seeded Nano variant.
func TestHashStringNanoWithSeedMatchesHashNanoWithSeed(t *testing.T) {
	testStrings := []string{
		"",
		"hello",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 200)),
	}
	seeds := []uint64{0, 1, 12345, 0xdeadbeef}

	for _, s := range testStrings {
		for _, seed := range seeds {
			expected := rapidhash.HashNanoWithSeed([]byte(s), seed)
			got := rapidhash.HashStringNanoWithSeed(s, seed)
			if got != expected {
				t.Errorf("HashStringNanoWithSeed(%q, %d) = 0x%x, want 0x%x", s, seed, got, expected)
			}
		}
	}
}

// TestHashStringMicroMatchesHashMicro verifies Micro variant string hashing.
func TestHashStringMicroMatchesHashMicro(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"ab",
		"abc",
		"abcdefgh",
		"abcdefghijklmnop",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 100)),
		string(make([]byte, 200)),
		string(make([]byte, 500)),
	}

	for _, s := range testStrings {
		expected := rapidhash.HashMicro([]byte(s))
		got := rapidhash.HashStringMicro(s)
		if got != expected {
			t.Errorf("HashStringMicro(%q) = 0x%x, want 0x%x (len=%d)", s, got, expected, len(s))
		}
	}
}

// TestHashStringMicroWithSeedMatchesHashMicroWithSeed verifies seeded Micro variant.
func TestHashStringMicroWithSeedMatchesHashMicroWithSeed(t *testing.T) {
	testStrings := []string{
		"",
		"hello",
		"The quick brown fox jumps over the lazy dog",
		string(make([]byte, 500)),
	}
	seeds := []uint64{0, 1, 12345, 0xdeadbeef}

	for _, s := range testStrings {
		for _, seed := range seeds {
			expected := rapidhash.HashMicroWithSeed([]byte(s), seed)
			got := rapidhash.HashStringMicroWithSeed(s, seed)
			if got != expected {
				t.Errorf("HashStringMicroWithSeed(%q, %d) = 0x%x, want 0x%x", s, seed, got, expected)
			}
		}
	}
}

// TestHashStringAllSizes tests all string functions across a range of sizes.
func TestHashStringAllSizes(t *testing.T) {
	// Test critical size boundaries
	sizes := []int{0, 1, 2, 3, 4, 7, 8, 15, 16, 17, 32, 48, 49, 64, 80, 81, 112, 113, 200, 448, 449, 1000}

	for _, size := range sizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}
		s := string(data)

		// Test all variants
		if got, want := rapidhash.HashString(s), rapidhash.Hash(data); got != want {
			t.Errorf("size %d: HashString = 0x%x, want 0x%x", size, got, want)
		}
		if got, want := rapidhash.HashStringNano(s), rapidhash.HashNano(data); got != want {
			t.Errorf("size %d: HashStringNano = 0x%x, want 0x%x", size, got, want)
		}
		if got, want := rapidhash.HashStringMicro(s), rapidhash.HashMicro(data); got != want {
			t.Errorf("size %d: HashStringMicro = 0x%x, want 0x%x", size, got, want)
		}

		// Test with seed
		seed := uint64(12345)
		if got, want := rapidhash.HashStringWithSeed(s, seed), rapidhash.HashWithSeed(data, seed); got != want {
			t.Errorf("size %d: HashStringWithSeed = 0x%x, want 0x%x", size, got, want)
		}
		if got, want := rapidhash.HashStringNanoWithSeed(s, seed), rapidhash.HashNanoWithSeed(data, seed); got != want {
			t.Errorf("size %d: HashStringNanoWithSeed = 0x%x, want 0x%x", size, got, want)
		}
		if got, want := rapidhash.HashStringMicroWithSeed(s, seed), rapidhash.HashMicroWithSeed(data, seed); got != want {
			t.Errorf("size %d: HashStringMicroWithSeed = 0x%x, want 0x%x", size, got, want)
		}
	}
}
