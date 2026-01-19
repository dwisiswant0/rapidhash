package rapidhash_test

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"reflect"
	"testing"

	"go.dw1.io/rapidhash"
)

type comparableStruct struct {
	A int
	B string
	C [3]uint16
}

func TestHashComparableSelfConsistent(t *testing.T) {
	seed := uint64(0x9e3779b97f4a7c15)

	x := 42
	arr := [4]uint32{1, 2, 3, 4}
	st := comparableStruct{A: 7, B: "hi", C: [3]uint16{9, 10, 11}}
	var nilPtr *int

	cases := []struct {
		name  string
		value any
	}{
		{"int", int(-12345)},
		{"uint64", uint64(0xdeadbeefcafebabe)},
		{"string-empty", ""},
		{"string", "rapidhash"},
		{"bool", true},
		{"uintptr", uintptr(123456)},
		{"ptr", &x},
		{"ptr-nil", nilPtr},
		{"array", arr},
		{"struct", st},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := hashComparableWithSeedAny(t, tc.value, seed)
			got2 := hashComparableWithSeedAny(t, tc.value, seed)
			if got1 != got2 {
				t.Fatalf("HashComparableWithSeed not deterministic: 0x%x vs 0x%x", got1, got2)
			}

			gotZero := hashComparableWithSeedAny(t, tc.value, 0)
			gotDefault := hashComparableAny(t, tc.value)
			if gotZero != gotDefault {
				t.Fatalf("HashComparable != HashComparableWithSeed(0): 0x%x vs 0x%x", gotDefault, gotZero)
			}
		})
	}
}

func TestHashComparableNaN(t *testing.T) {
	seed := uint64(123)
	v := math.NaN()

	h1 := rapidhash.HashComparableWithSeed(v, seed)
	h2 := rapidhash.HashComparableWithSeed(v, seed)

	// For NaN, equality is not required. This test only asserts no panic.
	_ = h1
	_ = h2
}

func TestWriteComparable(t *testing.T) {
	seed := uint64(0xfeedfacecafebeef)
	v := comparableStruct{A: 123, B: "abc", C: [3]uint16{1, 2, 3}}

	h := rapidhash.NewWithSeed(seed)
	h.WriteComparable(v)

	encoded := encodeComparableTest(v)
	expected := rapidhash.HashWithSeed(encoded, seed)

	if got := h.Sum64(); got != expected {
		t.Fatalf("WriteComparable hash = 0x%x, want 0x%x", got, expected)
	}
}

func TestWriteComparableWithOtherWrites(t *testing.T) {
	seed := uint64(0x1122334455667788)
	v := uint64(0x0102030405060708)
	prefix := []byte("prefix-")
	suffix := "-suffix"

	h := rapidhash.NewWithSeed(seed)
	_, _ = h.Write(prefix)
	h.WriteComparable(v)
	_, _ = h.WriteString(suffix)

	encoded := encodeComparableTest(v)
	combined := append(append(append([]byte{}, prefix...), encoded...), suffix...)
	expected := rapidhash.HashWithSeed(combined, seed)

	if got := h.Sum64(); got != expected {
		t.Fatalf("WriteComparable mixed hash = 0x%x, want 0x%x", got, expected)
	}
}

func hashComparableAny(t *testing.T, v any) uint64 {
	t.Helper()
	switch x := v.(type) {
	case int:
		return rapidhash.HashComparable(x)
	case uint64:
		return rapidhash.HashComparable(x)
	case string:
		return rapidhash.HashComparable(x)
	case bool:
		return rapidhash.HashComparable(x)
	case uintptr:
		return rapidhash.HashComparable(x)
	case *int:
		return rapidhash.HashComparable(x)
	case [4]uint32:
		return rapidhash.HashComparable(x)
	case comparableStruct:
		return rapidhash.HashComparable(x)
	default:
		t.Fatalf("unsupported test type %T", v)
		return 0
	}
}

func hashComparableWithSeedAny(t *testing.T, v any, seed uint64) uint64 {
	t.Helper()
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
		t.Fatalf("unsupported test type %T", v)
		return 0
	}
}

func encodeComparableTest(v any) []byte {
	return appendValueBytesTest(nil, reflect.ValueOf(v))
}

func appendValueBytesTest(buf []byte, v reflect.Value) []byte {
	buf = append(buf, v.Type().String()...)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return appendUint64LETest(buf, uint64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return appendUint64LETest(buf, v.Uint())
	case reflect.Array:
		var tmp [8]byte
		for i := 0; i < v.Len(); i++ {
			binary.LittleEndian.PutUint64(tmp[:], uint64(i))
			buf = append(buf, tmp[:]...)
			buf = appendValueBytesTest(buf, v.Index(i))
		}
		return buf
	case reflect.String:
		return append(buf, v.String()...)
	case reflect.Struct:
		var tmp [8]byte
		for i := 0; i < v.NumField(); i++ {
			binary.LittleEndian.PutUint64(tmp[:], uint64(i))
			buf = append(buf, tmp[:]...)
			buf = appendValueBytesTest(buf, v.Field(i))
		}
		return buf
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		buf = appendFloat64BytesTest(buf, real(c))
		buf = appendFloat64BytesTest(buf, imag(c))
		return buf
	case reflect.Float32, reflect.Float64:
		return appendFloat64BytesTest(buf, v.Float())
	case reflect.Bool:
		if v.Bool() {
			return append(buf, 1)
		}
		return append(buf, 0)
	case reflect.UnsafePointer, reflect.Pointer, reflect.Chan:
		return appendUint64LETest(buf, uint64(v.Pointer()))
	case reflect.Interface:
		return appendValueBytesTest(buf, v.Elem())
	}

	return buf
}

func appendFloat64BytesTest(buf []byte, f float64) []byte {
	if f == 0 {
		return append(buf, 0)
	}
	if math.IsNaN(f) {
		return appendUint64LETest(buf, randUint64Test())
	}
	return appendUint64LETest(buf, math.Float64bits(f))
}

func appendUint64LETest(buf []byte, x uint64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], x)
	return append(buf, tmp[:]...)
}

func randUint64Test() uint64 {
	var tmp [8]byte
	_, _ = rand.Read(tmp[:])
	return binary.LittleEndian.Uint64(tmp[:])
}
