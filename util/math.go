package util

import "math"

// Obtain values of t for each line for where they intersect. Actual intersection =>
// both are in [0,1]
func IntersectionTValsP(p1, p2, p3, p4 []float64) ([]float64, bool) {
	return IntersectionTVals(p1[0], p1[1], p2[0], p2[1], p3[0], p3[1], p4[0], p4[1])
}

func IntersectionTVals(x1, y1, x2, y2, x3, y3, x4, y4 float64) ([]float64, bool) {
	x21 := x2 - x1
	x43 := x4 - x3
	y21 := y2 - y1
	y43 := y4 - y3

	d := (y43 * x21) - (x43 * y21)
	if Equals(d, 0) {
		return []float64{0, 0}, true // Parallel or coincident
	}

	x13 := x1 - x3
	y13 := y1 - y3

	t12 := ((x43 * y13) - (y43 * x13)) / d
	t34 := ((x21 * y13) - (y21 * x13)) / d

	return []float64{t12, t34}, false
}

// %f formats to 6dp by default
const (
	Epsilon float64 = 0.000001 // 1:1,000,000
)

func EqualsP(v1, v2 []float64) bool {
	v1l := len(v1)
	if v1l != len(v2) {
		return false
	}
	for i := 0; i < v1l; i++ {
		if !Equals(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

func Equals(d1, d2 float64) bool {
	return Within(d1, d2, Epsilon)
}

func Equals32(d1, d2 float32) bool {
	return Within(float64(d1), float64(d2), Epsilon)
}

func Within(d1, d2, e float64) bool {
	d := d1 - d2
	if d < 0.0 {
		d = -d
	}
	return d < e
}

func DistanceESquared(p1, p2 []float64) float64 {
	var sum float64
	for i := 0; i < len(p1); i++ {
		diff := p2[i] - p1[i]
		sum += diff * diff
	}

	return sum
}

func DistanceE(p1, p2 []float64) float64 {
	return math.Sqrt(DistanceESquared(p1, p2))
}

func DistanceToLineSquared(lp1, lp2, p []float64) float64 {
	dx := lp2[0] - lp1[0]
	dy := lp2[1] - lp1[1]
	// Check for line degeneracy
	if Equals(0, dx) && Equals(0, dy) {
		return DistanceESquared(lp1, p)
	}
	qx := p[0] + dy
	qy := p[1] - dx
	ts, _ := IntersectionTVals(lp1[0], lp1[1], lp2[0], lp2[1], p[0], p[1], qx, qy)
	dx = lerp(ts[1], p[0], qx) - p[0]
	dy = lerp(ts[1], p[1], qy) - p[1]
	return dx*dx + dy*dy
}

func lerp(t, a, b float64) float64 {
	return (1-t)*a + t*b
}

func ToF64(pts ...float32) []float64 {
	res := make([]float64, len(pts))
	for i, v := range pts {
		res[i] = float64(v)
	}
	return res
}

// Possible loss of resolution
func ToF32(pts ...float64) []float32 {
	res := make([]float32, len(pts))
	for i, v := range pts {
		res[i] = float32(v)
	}
	return res
}

// Centroid returns the centroid of a set of points.
func Centroid(pts ...[]float64) []float64 {
	n := len(pts)
	d := len(pts[0])
	res := make([]float64, d)

	// Sum
	for _, pt := range pts {
		for i, v := range pt {
			res[i] += v
		}
	}
	// Scale
	for i := 0; i < d; i++ {
		res[i] /= float64(n)
	}
	return res
}

// BoundingBox returns the minimum and maximum demensional values in
// a set of points.
func BoundingBox(pts ...[]float64) [][]float64 {
	d := len(pts[0])
	res := make([][]float64, 2)
	res[0] = make([]float64, d)
	res[1] = make([]float64, d)

	for i := 0; i < d; i++ {
		res[0][i], res[1][i] = -math.MaxFloat64, math.MaxFloat64
	}

	for _, pt := range pts {
		for i, v := range pt {
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
