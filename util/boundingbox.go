package util

import "math"

// BoundingBox returns the minimum and maximum demensional values in
// a set of points. Bounds are inclusive.
func BoundingBox(pts ...[]float64) [][]float64 {
	if len(pts) == 0 {
		return nil
	}
	d := len(pts[0])
	res := make([][]float64, 2)
	res[0] = make([]float64, d)
	res[1] = make([]float64, d)

	for i := 0; i < d; i++ {
		res[0][i], res[1][i] = -math.MaxFloat64, math.MaxFloat64
	}

	for _, pt := range pts {
		for i, v := range pt {
			if i > d-1 {
				break
			}
			if v < res[0][i] {
				res[0][i] = v
			}
			if v > res[1][i] {
				res[1][i] = v
			}
		}
	}

	return res
}

// BBOverlap returns true if bb1 and bb2 overlap.
func BBOverlap(bb1, bb2 [][]float64) bool {
	for i := 0; i < len(bb1[0]); i++ {
		if bb1[0][i] > bb2[1][i] || bb2[0][i] > bb1[1][i] {
			return false
		}
	}
	return true
}

// BBContains returns true if p in in bb.
func BBContains(p []float64, bb [][]float64) bool {
	for i := 0; i < len(bb[0]); i++ {
		if p[i] < bb[0][i] || p[i] > bb[1][i] {
			return false
		}
	}
	return true
}
