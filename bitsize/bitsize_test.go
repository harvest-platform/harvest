package bitsize

import (
	"math"
	"testing"
)

var floatLimits = map[string][2]float64{
	"maxfloat32": {math.MaxFloat32, 32},
	"minfloat32": {math.SmallestNonzeroFloat32, 32},
	"maxfloat64": {math.MaxFloat64, 64},
	"minfloat64": {math.SmallestNonzeroFloat64, 64},
}

var intLimits = map[string][2]int64{
	"maxint2":  {1<<1 - 1, 2},
	"minint2":  {-1 << 1, 2},
	"maxint4":  {1<<3 - 1, 4},
	"minint4":  {-1 << 3, 4},
	"maxint8":  {math.MaxInt8, 8},
	"minint8":  {math.MinInt8, 8},
	"maxint16": {math.MaxInt16, 16},
	"minint16": {math.MinInt16, 16},
	"maxint32": {math.MaxInt32, 32},
	"minint32": {math.MinInt32, 32},
	"maxint64": {math.MaxInt64, 64},
	"minint64": {math.MinInt64, 64},
}

var uintLimits = map[string][2]uint64{
	"maxuint2":  {1<<2 - 1, 2},
	"maxuint4":  {1<<4 - 1, 4},
	"maxuint8":  {math.MaxUint8, 8},
	"maxuint16": {math.MaxUint16, 16},
	"maxuint32": {math.MaxUint32, 32},
	"maxuint64": {math.MaxUint64, 64},
}

func TestFloat(t *testing.T) {
	for name, limit := range floatLimits {
		t.Run(name, func(t *testing.T) {
			act := Float(limit[0])
			exp := uint8(limit[1])

			if act != exp {
				t.Errorf("expected %d, got %d", exp, act)
			}
		})
	}
}

func TestInt(t *testing.T) {
	for name, limit := range intLimits {
		t.Run(name, func(t *testing.T) {
			act := Int(limit[0])
			exp := uint8(limit[1])

			if act != exp {
				t.Errorf("expected %d, got %d", exp, act)
			}
		})
	}
}

func TestUint(t *testing.T) {
	for name, limit := range uintLimits {
		t.Run(name, func(t *testing.T) {
			act := Uint(limit[0])
			exp := uint8(limit[1])

			if act != exp {
				t.Errorf("expected %d, got %d", exp, act)
			}
		})
	}
}

func BenchmarkFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float(-32932.92329032930)
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int(-32382897329932)
	}
}

func BenchmarkUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint(32382897329932)
	}
}
