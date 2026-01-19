// Package rapidhash provides a Go implementation of the rapidhash V3 algorithm.
//
// Rapidhash is a very fast, high quality, platform-independent hashing algorithm.
// It is based on 'wyhash' by Wang Yi.
//
// This package provides three hash variants optimized for different use cases:
//
//   - [Hash]/[HashWithSeed]: The default variant, optimal for large inputs
//     (>112 bytes). Uses 7 parallel mixing lanes, processing 112 bytes per
//     iteration.
//
//   - [HashMicro]/[HashMicroWithSeed]: Optimized for cache-sensitive HPC/server
//     cases. Uses 5 parallel lanes and is typically fastest up to ~512 bytes;
//     above that it is usually within a few percent of [Hash] but can be slower
//     depending on CPU and toolchain. Prefer [Hash] for consistently large inputs.
//
//   - [HashNano]/[HashNanoWithSeed]: Optimized for mobile/embedded with minimal
//     code size. Uses 3 parallel lanes, fastest for inputs up to 48 bytes.
//
// # Variant Compatibility
//
// All three variants produce identical output for small inputs:
//   - All match for inputs ≤ 48 bytes.
//   - [Hash] and [HashMicro] match for inputs ≤ 80 bytes.
//
// For larger inputs, each variant produces different (but equally valid) hashes.
//
// The comparable hashing helpers ([HashComparable], [HashComparableWithSeed],
// and [Hasher.WriteComparable]) use a different encoding strategy (type tagging
// plus reflection traversal), so their outputs are not compatible with [Hash]
// or [HashWithSeed]. They also randomize floating-point NaNs and hash
// pointer-like values by address, which can make results non-deterministic or
// process-specific.
//
// # Performance
//
// On modern x86-64 CPUs, typical performance is:
//   - Small inputs (8-16 bytes): ~2-4 GiB/s.
//   - Medium inputs (64-256 bytes): ~8-13.5 GiB/s.
//   - Large inputs (1KB+): ~15-21 GiB/s.
//
// Performance is hardware- and toolchain-dependent; your results may vary on
// different CPUs, microarchitectures, and Go versions.
//
// # Thread Safety
//
// All hash functions ([Hash], [HashMicro], [HashNano], etc.) are safe for
// concurrent use. The [Hasher] type is NOT safe for concurrent use; each
// goroutine should use its own [Hasher] instance.
//
// # References
//
//   - Original rapidhash: https://github.com/Nicoshev/rapidhash
//   - Based on wyhash: https://github.com/wangyi-fudan/wyhash
package rapidhash
