package util

import "math"

// BoundingBox returns the minimum and maximum dimensional values in
// a set of points. Bounds are inclusive. The dimensionality of the
// smallest dimension point is used for all points.
func BoundingBox(pts ...[]float64) [][]float64 {
	if len(pts) == 0 {
		return nil
	}

	d := MinD(pts...)
	res := make([][]float64, 2)
	res[0] = make([]float64, d)
	res[1] = make([]float64, d)

	for i := 0; i < d; i++ {
		res[0][i], res[1][i] = math.MaxFloat64, -math.MaxFloat64
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

// BBOverlap returns true if bb1 and bb2 overlap at the smallest dimensionality.
func BBOverlap(bb1, bb2 [][]float64) bool {
	min := len(bb1[0])
	if min > len(bb2[0]) {
		min = len(bb2[0])
	}
	for i := 0; i < min; i++ {
		if bb1[0][i] > bb2[1][i] || bb2[0][i] > bb1[1][i] {
			return false
		}
	}
	return true
}

// BBContains returns true if p is in bb at the smallest dimensionality.
func BBContains(p []float64, bb [][]float64) bool {
	min := len(bb[0])
	if min > len(p) {
		min = len(p)
	}
	for i := 0; i < min; i++ {
		if p[i] < bb[0][i] || p[i] > bb[1][i] {
			return false
		}
	}
	return true
}
