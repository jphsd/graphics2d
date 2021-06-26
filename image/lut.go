package image

import (
	"fmt"
	"github.com/jphsd/graphics2d/util"
)

// CreateLutFromValues maps a series of values in [0,1] to [0,255]
// Note - no checking on values range
func CreateLutFromValues(values []float64) []uint8 {
	n := len(values)
	res := make([]uint8, n)
	for i := 0; i < n; i++ {
		res[i] = uint8(values[i] * 255)
	}
	return res
}

// PremulLut creates a slice of luts from lut for every possible value of alpha
func PremulLut(lut []uint8) [][]uint8 {
	n := len(lut)
	res := make([][]uint8, n)
	for i := 0; i < n; i++ {
		nres := make([]uint8, 256)
		v := uint32(lut[i])
		v |= v << 8
		for a := 0; a < 256; a++ {
			nv := uint32(a) * v
			nv /= 0xff
			nres[a] = uint8(nv >> 8)
		}
		res[i] = nres
	}
	return res
}

// NLExpansionLut generates a lut [start,end), normalized and mapped through f.
// For example NLExpansionLut(256, 0, 256, &NLSin{}) will generate a Sin ramp.
func NLExpansionLut(n, start, end int, f util.NonLinear) []uint8 {
	if start > end {
		start, end = end, start
	}
	if start < 0 || start > n-1 || end < 1 || end > n {
		panic(fmt.Errorf("start or end not in range"))
	}
	res := make([]uint8, n)

	var i int
	for ; i < start; i++ {
		res[i] = 0
	}
	delta := end - start
	if delta > 0 {
		div := 1 / float64(delta)
		for ; i < end; i++ {
			v := float64(i-start) * div
			v = f.Transform(v)
			res[i] = uint8(v * 0xff)
		}
	}
	for ; i < n; i++ {
		res[i] = 0xff
	}

	return res
}
