# rapidhash

[![Go Reference](https://pkg.go.dev/badge/go.dw1.io/rapidhash.svg)](https://pkg.go.dev/go.dw1.io/rapidhash)

a Go implementation of the rapidhash V3 algorithm - a very fast, high quality, platform-independent hashing algorithm.

## Install

```bash
go get go.dw1.io/rapidhash
```

## Usage

### Basic Hashing

```go
package main

import (
    "fmt"
    "go.dw1.io/rapidhash"
)

func main() {
    data := []byte("hello world")

    // default seed (0)
    hash := rapidhash.Hash(data)
    fmt.Printf("Hash: 0x%x\n", hash)

    // custom seed
    hash = rapidhash.HashWithSeed(data, 12345)
    fmt.Printf("Hash with seed: 0x%x\n", hash)
}
```

### Variant Selection

```go
// for small inputs (≤48 bytes) - fastest for mobile/embedded
nano := rapidhash.HashNano([]byte("key"))
fmt.Printf("Nano: 0x%x\n", nano)

// for medium inputs (≤512 bytes) - optimized for HPC/server
micro := rapidhash.HashMicro([]byte("medium data"))
fmt.Printf("Micro: 0x%x\n", micro)

// for large inputs (>512 bytes) - general purpose
large := rapidhash.Hash([]byte("large input data..."))
fmt.Printf("Large: 0x%x\n", large)
```

### Streaming Hash

```go
// for incremental hashing
hasher := rapidhash.New()
hasher.Write([]byte("hello "))
hasher.Write([]byte("world"))
hash := hasher.Sum64()
fmt.Printf("Streaming hash: 0x%x\n", hash)

// reset and reuse
hasher.Reset()
hasher.Write([]byte("new data"))
hash = hasher.Sum64()
```

## Performance

Typical performance on modern x86-64 CPUs (AMD EPYC 7763):

* Small inputs (8-16 bytes): **~2-4 GiB/s**.
* Medium inputs (64-256 bytes): **~8-13.5 GiB/s**.
* Large inputs (1KB+): **~15-21 GiB/s**.

Performance varies by hardware, microarchitecture, and Go version.

## Benchmarks

<details open>
  <summary><code>benchstat</code></summary>

  ```
  goos: linux
  goarch: amd64
  pkg: go.dw1.io/rapidhash
  cpu: AMD EPYC 7763 64-Core Processor                
          │     Hash     │               HashNano                │               HashMicro               │                Hasher                 │
          │    sec/op    │   sec/op     vs base                  │   sec/op     vs base                  │    sec/op     vs base                 │
  8-4        3.741n ± 0%   3.740n ± 0%        ~ (p=0.559 n=10)                                              9.677n ± 0%  +158.67% (p=0.000 n=10)
  16-4       3.743n ± 0%   3.741n ± 0%        ~ (p=0.236 n=10)                                              9.976n ± 0%  +166.57% (p=0.000 n=10)
  32-4       7.488n ± 0%   4.984n ± 0%  -33.44% (p=0.000 n=10)                                             14.040n ± 0%   +87.50% (p=0.000 n=10)
  64-4       7.788n ± 0%                                           5.611n ± 0%  -27.96% (p=0.000 n=10)     14.675n ± 2%   +88.42% (p=0.000 n=10)
  128-4     11.540n ± 1%                                           9.670n ± 0%  -16.20% (p=0.000 n=10)     19.550n ± 0%   +69.41% (p=0.000 n=10)
  256-4      17.77n ± 0%                                           15.77n ± 0%  -11.28% (p=0.000 n=10)      27.64n ± 0%   +55.47% (p=0.000 n=10)
  512-4      30.05n ± 0%                                           28.53n ± 1%   -5.06% (p=0.000 n=10)      41.16n ± 0%   +36.99% (p=0.000 n=10)
  1024-4     52.05n ± 0%                                                                                    68.17n ± 0%   +30.97% (p=0.000 n=10)
  4096-4     185.9n ± 0%                                                                                    231.9n ± 0%   +24.72% (p=0.000 n=10)
  8192-4     365.4n ± 0%                                                                                    451.1n ± 0%   +23.48% (p=0.000 n=10)
  geomean    21.14n        4.116n       -12.71%                ¹   12.50n       -15.56%                ¹    35.47n        +67.83%
  ¹ benchmark set differs from baseline; geomeans may not be comparable

          │     Hash      │                HashNano                 │                HashMicro                │                Hasher                │
          │      B/s      │      B/s       vs base                  │      B/s       vs base                  │     B/s       vs base                │
  8-4       2039.3Mi ± 0%   2040.1Mi ± 0%        ~ (p=0.579 n=10)                                               788.4Mi ± 0%  -61.34% (p=0.000 n=10)
  16-4       3.981Gi ± 0%    3.983Gi ± 0%        ~ (p=0.239 n=10)                                               1.494Gi ± 0%  -62.48% (p=0.000 n=10)
  32-4       3.980Gi ± 0%    5.979Gi ± 0%  +50.23% (p=0.000 n=10)                                               2.123Gi ± 0%  -46.66% (p=0.000 n=10)
  64-4       7.653Gi ± 0%                                             10.623Gi ± 0%  +38.81% (p=0.000 n=10)     4.062Gi ± 2%  -46.92% (p=0.000 n=10)
  128-4     10.333Gi ± 1%                                             12.328Gi ± 0%  +19.31% (p=0.000 n=10)     6.098Gi ± 0%  -40.98% (p=0.000 n=10)
  256-4     13.413Gi ± 0%                                             15.119Gi ± 0%  +12.72% (p=0.000 n=10)     8.628Gi ± 0%  -35.67% (p=0.000 n=10)
  512-4      15.87Gi ± 0%                                              16.72Gi ± 1%   +5.33% (p=0.000 n=10)     11.58Gi ± 0%  -27.01% (p=0.000 n=10)
  1024-4     18.32Gi ± 0%                                                                                       13.99Gi ± 0%  -23.65% (p=0.000 n=10)
  4096-4     20.52Gi ± 0%                                                                                       16.45Gi ± 0%  -19.82% (p=0.000 n=10)
  8192-4     20.88Gi ± 0%                                                                                       16.91Gi ± 0%  -19.02% (p=0.000 n=10)
  geomean    9.163Gi         3.620Gi       +14.56%                ¹    13.49Gi       +18.41%                ¹   5.460Gi       -40.42%
  ¹ benchmark set differs from baseline; geomeans may not be comparable
  ```

  > All bytes and allocs per operation samples are all zero and equal.
</details>

The [`HashNano`](https://pkg.go.dev/go.dw1.io/rapidhash#HashNano) variant is optimized for small inputs (≤48 bytes), [`HashMicro`](https://pkg.go.dev/go.dw1.io/rapidhash#HashMicro) for medium inputs (≤512 bytes), and [`Hash`](https://pkg.go.dev/go.dw1.io/rapidhash#Hash) for general use. The streaming [`Hasher`](https://pkg.go.dev/go.dw1.io/rapidhash#Hasher) is slightly slower but provides incremental hashing capabilities.

### Comparison

<details open>
  <summary><code>benchstat</code></summary>

  ```
  goos: linux
  goarch: amd64
  pkg: benchmarks
  cpu: AMD EPYC 7763 64-Core Processor                
          │  rapidhash  │               wyhash                │                xxh64                 │
          │   sec/op    │   sec/op     vs base                │    sec/op     vs base                │
  8-4       3.734n ± 0%   3.737n ± 0%        ~ (p=0.209 n=10)    4.679n ± 0%  +25.31% (p=0.000 n=10)
  16-4      3.737n ± 0%   4.359n ± 0%  +16.63% (p=0.000 n=10)    5.298n ± 0%  +41.77% (p=0.000 n=10)
  32-4      7.167n ± 0%   5.604n ± 0%  -21.80% (p=0.000 n=10)    9.601n ± 0%  +33.98% (p=0.000 n=10)
  64-4      7.476n ± 0%   9.027n ± 0%  +20.74% (p=0.000 n=10)   12.150n ± 0%  +62.52% (p=0.000 n=10)
  128-4     11.21n ± 0%   12.77n ± 0%  +13.92% (p=0.000 n=10)    17.13n ± 0%  +52.81% (p=0.000 n=10)
  256-4     17.45n ± 0%   20.55n ± 0%  +17.77% (p=0.000 n=10)    27.08n ± 0%  +55.19% (p=0.000 n=10)
  512-4     30.15n ± 0%   35.12n ± 1%  +16.48% (p=0.000 n=10)    47.00n ± 0%  +55.87% (p=0.000 n=10)
  1024-4    51.62n ± 0%   64.38n ± 0%  +24.70% (p=0.000 n=10)    86.86n ± 0%  +68.25% (p=0.000 n=10)
  4096-4    186.2n ± 0%   243.0n ± 0%  +30.48% (p=0.000 n=10)    328.5n ± 0%  +76.42% (p=0.000 n=10)
  8192-4    365.0n ± 0%   479.2n ± 0%  +31.27% (p=0.000 n=10)    646.7n ± 0%  +77.14% (p=0.000 n=10)
  geomean   20.84n        23.74n       +13.91%                   32.10n       +54.03%

          │   rapidhash   │                wyhash                 │                xxh64                 │
          │      B/s      │      B/s       vs base                │     B/s       vs base                │
  8-4        1.995Gi ± 0%    1.994Gi ± 0%        ~ (p=0.271 n=10)   1.592Gi ± 0%  -20.19% (p=0.000 n=10)
  16-4       3.988Gi ± 0%    3.419Gi ± 0%  -14.26% (p=0.000 n=10)   2.812Gi ± 0%  -29.47% (p=0.000 n=10)
  32-4       4.158Gi ± 0%    5.318Gi ± 0%  +27.88% (p=0.000 n=10)   3.104Gi ± 0%  -25.36% (p=0.000 n=10)
  64-4       7.973Gi ± 0%    6.603Gi ± 0%  -17.17% (p=0.000 n=10)   4.905Gi ± 0%  -38.48% (p=0.000 n=10)
  128-4     10.636Gi ± 0%    9.336Gi ± 0%  -12.23% (p=0.000 n=10)   6.958Gi ± 0%  -34.58% (p=0.000 n=10)
  256-4     13.666Gi ± 0%   11.603Gi ± 0%  -15.10% (p=0.000 n=10)   8.805Gi ± 0%  -35.57% (p=0.000 n=10)
  512-4      15.82Gi ± 0%    13.58Gi ± 1%  -14.16% (p=0.000 n=10)   10.15Gi ± 0%  -35.85% (p=0.000 n=10)
  1024-4     18.47Gi ± 0%    14.81Gi ± 0%  -19.81% (p=0.000 n=10)   10.98Gi ± 0%  -40.57% (p=0.000 n=10)
  4096-4     20.49Gi ± 0%    15.70Gi ± 0%  -23.35% (p=0.000 n=10)   11.61Gi ± 0%  -43.32% (p=0.000 n=10)
  8192-4     20.90Gi ± 0%    15.92Gi ± 0%  -23.82% (p=0.000 n=10)   11.80Gi ± 0%  -43.55% (p=0.000 n=10)
  geomean    9.292Gi         8.157Gi       -12.21%                  6.032Gi       -35.08%
  ```

  > All bytes and allocs per operation samples are all zero and equal.
</details>

**rapidhash** outperforms wyhash by **~14%** and xxh64 by **~55%** on average across various input sizes. It excels particularly for larger inputs, achieving up to **~21 GiB/s**.

Run benchmarks yourself:

```bash
make bench
make bench -C benchmarks
```

### Comparable Hash Comparison

<details open>
  <summary><code>benchstat</code></summary>

  ```
  goos: linux
  goarch: amd64
  pkg: benchmarks
  cpu: AMD EPYC 7763 64-Core Processor                
                  │ rapidhash.HashComparable │         maphash.Comparable          │
                  │          sec/op          │   sec/op     vs base                │
  Comparable/int-4                    33.50n ± 3%   43.37n ± 1%  +29.45% (p=0.000 n=10)
  Comparable/uint64-4                 34.63n ± 0%   44.15n ± 1%  +27.49% (p=0.000 n=10)
  Comparable/string-4                 51.58n ± 1%   58.87n ± 0%  +14.14% (p=0.000 n=10)
  Comparable/bool-4                   23.70n ± 0%   27.15n ± 0%  +14.58% (p=0.000 n=10)
  Comparable/uintptr-4                33.88n ± 0%   44.20n ± 2%  +30.44% (p=0.000 n=10)
  Comparable/ptr-4                    24.88n ± 0%   33.35n ± 2%  +34.07% (p=0.000 n=10)
  Comparable/ptr-nil-4                24.87n ± 0%   33.35n ± 0%  +34.10% (p=0.000 n=10)
  Comparable/array-4                  124.3n ± 0%   162.9n ± 0%  +31.05% (p=0.000 n=10)
  Comparable/struct-4                 181.1n ± 0%   252.8n ± 0%  +39.56% (p=0.000 n=10)
  geomean                             44.39n        56.84n       +28.05%
                  │ rapidhash.HashComparable │         maphash.Comparable          │
                  │           B/op           │    B/op     vs base                 │
  Comparable/int-4                   8.000 ± 0%     8.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/uint64-4                8.000 ± 0%     8.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/string-4                16.00 ± 0%     16.00 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/bool-4                  0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/uintptr-4               8.000 ± 0%     8.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/ptr-4                   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/ptr-nil-4               0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/array-4                 16.00 ± 0%     16.00 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/struct-4                32.00 ± 0%     32.00 ± 0%       ~ (p=1.000 n=10) ¹
  geomean                                       ²               +0.00%                ²
  ¹ all samples are equal
  ² summaries must be >0 to compute geomean
                  │ rapidhash.HashComparable │         maphash.Comparable          │
                  │        allocs/op         │ allocs/op   vs base                 │
  Comparable/int-4                   1.000 ± 0%     1.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/uint64-4                1.000 ± 0%     1.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/string-4                1.000 ± 0%     1.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/bool-4                  0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/uintptr-4               1.000 ± 0%     1.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/ptr-4                   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/ptr-nil-4               0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/array-4                 1.000 ± 0%     1.000 ± 0%       ~ (p=1.000 n=10) ¹
  Comparable/struct-4                1.000 ± 0%     1.000 ± 0%       ~ (p=1.000 n=10) ¹
  geomean                                       ²               +0.00%                ²
  ¹ all samples are equal
  ² summaries must be >0 to compute geomean
  ```
</details>

Across the benchmarked comparable types, **rapidhash** outperforms [`maphash`](https://pkg.go.dev/hash/maphash) in speed while matching its alloc behavior. Performance improvements differ by type, with notable gains for pointer-based and composite structures. Note that the std lib's [`maphash`](https://pkg.go.dev/hash/maphash) was tested using the `-tags=purego` flag, which uses the wyhash algorithm.

Run benchmarks yourself:

```bash
make bench-comparable -C benchmarks
```

## Variant Compatibility

All three variants produce identical output for small inputs:
* All match for inputs ≤ 48 bytes.
* [`Hash`](https://pkg.go.dev/go.dw1.io/rapidhash#Hash) and [`HashMicro`](https://pkg.go.dev/go.dw1.io/rapidhash#HashMicro) match for inputs ≤ 80 bytes.

The comparable hashing helpers ([`HashComparable`](https://pkg.go.dev/go.dw1.io/rapidhash#HashComparable),
[`HashComparableWithSeed`](https://pkg.go.dev/go.dw1.io/rapidhash#HashComparableWithSeed), and
[`Hasher.WriteComparable`](https://pkg.go.dev/go.dw1.io/rapidhash#Hasher.WriteComparable)) use a
different encoding strategy (type tagging plus reflection traversal), so their outputs are not
compatible with [`Hash`](https://pkg.go.dev/go.dw1.io/rapidhash#Hash) or
[`HashWithSeed`](https://pkg.go.dev/go.dw1.io/rapidhash#HashWithSeed). They also randomize
floating-point NaNs and hash pointer-like values by address, which can make results non-deterministic
or process-specific.

## Thread Safety

* All hash functions ([`Hash`](https://pkg.go.dev/go.dw1.io/rapidhash#Hash), [`HashMicro`](https://pkg.go.dev/go.dw1.io/rapidhash#HashMicro), [`HashNano`](https://pkg.go.dev/go.dw1.io/rapidhash#HashNano), etc.) are safe for concurrent use.
* [`Hasher`](https://pkg.go.dev/go.dw1.io/rapidhash#Hasher) instances are **not** safe for concurrent use - use one per goroutine.

## References

* **Original C Implementation**: [Nicoshev/rapidhash](https://github.com/Nicoshev/rapidhash).
* **Based on**: [wangyi-fudan/wyhash](https://github.com/wangyi-fudan/wyhash).

## License

MIT. See [LICENSE](/LICENSE).

