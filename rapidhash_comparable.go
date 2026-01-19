package rapidhash

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

// HashComparable returns the hash of comparable value v using the default seed (0).
//
// This is not compatible with [Hash] or [HashWithSeed] because it encodes type
// information and traverses values via reflection; it also randomizes
// floating-point NaNs (so results are not deterministic when v contains NaNs)
// and hashes pointer-like values by address, making results process-specific.
func HashComparable[T comparable](v T) uint64 {
	return HashComparableWithSeed(v, 0)
}

// HashComparableWithSeed returns the hash of comparable value v using seed.
//
// This is not compatible with [Hash] or [HashWithSeed] because it encodes type
// information and traverses values via reflection; it also randomizes
// floating-point NaNs (so results are not deterministic when v contains NaNs)
// and hashes pointer-like values by address, making results process-specific.
func HashComparableWithSeed[T comparable](v T, seed uint64) uint64 {
	var stack [256]byte
	buf := stack[:0]
	buf = appendComparableBytes(buf, reflect.ValueOf(v))

	return HashWithSeed(buf, seed)
}

func appendComparableBytes(buf []byte, v reflect.Value) []byte {
	return appendValueBytes(buf, v)
}

func appendValueBytes(buf []byte, v reflect.Value) []byte {
	buf = append(buf, v.Type().String()...)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return appendUint64LE(buf, uint64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return appendUint64LE(buf, v.Uint())
	case reflect.Array:
		var tmp [8]byte
		for i := 0; i < v.Len(); i++ {
			binary.LittleEndian.PutUint64(tmp[:], uint64(i))
			buf = append(buf, tmp[:]...)
			buf = appendValueBytes(buf, v.Index(i))
		}
		return buf
	case reflect.String:
		return append(buf, v.String()...)
	case reflect.Struct:
		var tmp [8]byte
		for i := 0; i < v.NumField(); i++ {
			binary.LittleEndian.PutUint64(tmp[:], uint64(i))
			buf = append(buf, tmp[:]...)
			buf = appendValueBytes(buf, v.Field(i))
		}
		return buf
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		buf = appendFloat64Bytes(buf, real(c))
		buf = appendFloat64Bytes(buf, imag(c))
		return buf
	case reflect.Float32, reflect.Float64:
		return appendFloat64Bytes(buf, v.Float())
	case reflect.Bool:
		if v.Bool() {
			return append(buf, 1)
		}
		return append(buf, 0)
	case reflect.UnsafePointer, reflect.Pointer, reflect.Chan:
		return appendUint64LE(buf, uint64(v.Pointer()))
	case reflect.Interface:
		return appendValueBytes(buf, v.Elem())
	}

	panic(errors.New("rapidhash: hash of unhashable type " + v.Type().String()))
}

func appendFloat64Bytes(buf []byte, f float64) []byte {
	if f == 0 {
		return append(buf, 0)
	}
	if math.IsNaN(f) {
		return appendUint64LE(buf, randUint64())
	}
	return appendUint64LE(buf, math.Float64bits(f))
}

func appendUint64LE(buf []byte, x uint64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], x)
	return append(buf, tmp[:]...)
}

func randUint64() uint64 {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	return binary.LittleEndian.Uint64(buf[:])
}
