package benchmarks_test

import (
	"hash/maphash"
	"testing"

	"go.dw1.io/rapidhash"
)

type comparableStruct struct {
	A int
	B string
	C [3]uint16
}

func BenchmarkComparable(b *testing.B) {
	seed := uint64(0x9e3779b97f4a7c15)
	mapSeed := maphash.MakeSeed()
	x := 42
	arr := [4]uint32{1, 2, 3, 4}
	st := comparableStruct{A: 7, B: "hi", C: [3]uint16{9, 10, 11}}
	var nilPtr *int

	cases := []struct {
		name  string
		value any
	}{
		{name: "int", value: int(-12345)},
		{name: "uint64", value: uint64(0xdeadbeefcafebabe)},
		{name: "string", value: "rapidhash"},
		{name: "bool", value: true},
		{name: "uintptr", value: uintptr(123456)},
		{name: "ptr", value: &x},
		{name: "ptr-nil", value: nilPtr},
		{name: "array", value: arr},
		{name: "struct", value: st},
	}

	for _, tc := range cases {
		b.Run("rapidhash/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sink = hashComparableWithSeedAnyBench(tc.value, seed)
			}
		})

		b.Run("maphash/"+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sink = maphashComparableAnyBench(tc.value, mapSeed)
			}
		})
	}
}

func hashComparableWithSeedAnyBench(v any, seed uint64) uint64 {
	switch x := v.(type) {
	case int:
		return rapidhash.HashComparableWithSeed(x, seed)
	case uint64:
		return rapidhash.HashComparableWithSeed(x, seed)
	case string:
		return rapidhash.HashComparableWithSeed(x, seed)
	case bool:
		return rapidhash.HashComparableWithSeed(x, seed)
	case uintptr:
		return rapidhash.HashComparableWithSeed(x, seed)
	case *int:
		return rapidhash.HashComparableWithSeed(x, seed)
	case [4]uint32:
		return rapidhash.HashComparableWithSeed(x, seed)
	case comparableStruct:
		return rapidhash.HashComparableWithSeed(x, seed)
	default:
		panic("unsupported test type")
	}
}

func maphashComparableAnyBench(v any, seed maphash.Seed) uint64 {
	switch x := v.(type) {
	case int:
		return maphash.Comparable(seed, x)
	case uint64:
		return maphash.Comparable(seed, x)
	case string:
		return maphash.Comparable(seed, x)
	case bool:
		return maphash.Comparable(seed, x)
	case uintptr:
		return maphash.Comparable(seed, x)
	case *int:
		return maphash.Comparable(seed, x)
	case [4]uint32:
		return maphash.Comparable(seed, x)
	case comparableStruct:
		return maphash.Comparable(seed, x)
	default:
		panic("unsupported test type")
	}
}
