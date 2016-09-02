package bitsize

import "math"

var sizes [65]uint64

func init() {
	for i := 0; i < 65; i++ {
		sizes[i] = 1<<uint8(i) - 1
	}
}

func Float(v interface{}) uint8 {
	switch f := v.(type) {
	case float32:
		return 32

	case float64:
		if f >= math.SmallestNonzeroFloat32 && f <= math.MaxFloat32 {
			return 32
		}

		return 64
	}

	panic("float required")
}

func Uint(x uint64) uint8 {
	for b, max := range sizes {
		if x <= max {
			return uint8(b)
		}
	}

	panic("unknown unsigned bitsize")
}

func Int(x int64) uint8 {
	if x < 0 {
		y := uint64(-x)

		for b, max := range sizes {
			if y <= (max + 1) {
				return uint8(b) + 1
			}
		}

		panic("unknown bitsize")
	}

	y := uint64(x)

	for b, max := range sizes {
		if y <= max {
			return uint8(b) + 1
		}
	}

	panic("unknown bitsize")
}
