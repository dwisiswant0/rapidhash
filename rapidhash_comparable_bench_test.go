package rapidhash_test

import (
	"testing"

	"go.dw1.io/rapidhash"
)

func BenchmarkComparable(b *testing.B) {
	seed := uint64(0x9e3779b97f4a7c15)
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
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sink = hashComparableWithSeedAnyBench(tc.value, seed)
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
