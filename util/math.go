package util

import (
	"fmt"
	"math"
)

// IntersectionTValsP obtains values of t for each line for where they intersect. Actual intersection =>
// both are in [0,1]
func IntersectionTValsP(p1, p2, p3, p4 []float64) ([]float64, error) {
	return IntersectionTVals(p1[0], p1[1], p2[0], p2[1], p3[0], p3[1], p4[0], p4[1])
}

// IntersectionTVals obtains values of t for each line for where they intersect. Actual intersection =>
// both are in [0,1]
func IntersectionTVals(x1, y1, x2, y2, x3, y3, x4, y4 float64) ([]float64, error) {
	x21 := x2 - x1
	x43 := x4 - x3
	y21 := y2 - y1
	y43 := y4 - y3

	d := (y43 * x21) - (x43 * y21)
	if Equals(d, 0) {
		return nil, fmt.Errorf("parallel or coincident")
	}

	x13 := x1 - x3
	y13 := y1 - y3

	t12 := ((x43 * y13) - (y43 * x13)) / d
	t34 := ((x21 * y13) - (y21 * x13)) / d

	return []float64{t12, t34}, nil
}

// %f formats to 6dp by default
const (
	Epsilon float64 = 0.000001 // 1:1,000,000
)

// EqualsP returns true if two points are equal.
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

// Equals returns true if two values are within Epsilon of each other.
func Equals(d1, d2 float64) bool {
	return Within(d1, d2, Epsilon)
}

// Equals32 is the float32 version of Equals.
func Equals32(d1, d2 float32) bool {
	return Within(float64(d1), float64(d2), Epsilon)
}

// Within returns true if the two values are within e of each other.
func Within(d1, d2, e float64) bool {
	d := d1 - d2
	if d < 0.0 {
		d = -d
	}
	return d < e
}

// DistanceESquared returns the squared Euclidean distance between two points.
func DistanceESquared(p1, p2 []float64) float64 {
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	return dx*dx + dy*dy
}

// DistanceESquaredN returns the squared Euclidean distance between two points.
func DistanceESquaredN(p1, p2 []float64) float64 {
	var sum float64
	for i := 0; i < MinD(p1, p2); i++ {
		diff := p2[i] - p1[i]
		sum += diff * diff
	}

	return sum
}

// DistanceE returns the Euclidean distance between two points.
func DistanceE(p1, p2 []float64) float64 {
	return math.Sqrt(DistanceESquared(p1, p2))
}

// DistanceToLineSquared calculates the squared Euclidean length of the normal from a point to
// the line.
func DistanceToLineSquared(lp1, lp2, p []float64) float64 {
	dx := lp2[0] - lp1[0]
	dy := lp2[1] - lp1[1]
	// Check for line degeneracy
	if Equals(0, dx) && Equals(0, dy) {
		return DistanceESquared(lp1, p)
	}
	qx := p[0] + dy
	qy := p[1] - dx
	ts, err := IntersectionTVals(lp1[0], lp1[1], lp2[0], lp2[1], p[0], p[1], qx, qy)
	if err != nil {
		return DistanceESquared(lp1, p)
	}
	dx = Lerp(ts[1], p[0], qx) - p[0]
	dy = Lerp(ts[1], p[1], qy) - p[1]
	return dx*dx + dy*dy
}

// SideOfLine calculates which side of a line a point is one by calculating the dot product of the
// vector from the line start to the point with the line's normal. If +ve then one side, -ve the other,
// 0 - on the line.
func SideOfLine(lp1, lp2, p []float64) float64 {
	return (p[0]-lp1[0])*(lp2[1]-lp1[1]) - (p[1]-lp1[1])*(lp2[0]-lp1[0])
}

// ToF64 casts a slice of float32 to float64.
func ToF64(pts ...float32) []float64 {
	res := make([]float64, len(pts))
	for i, v := range pts {
		res[i] = float64(v)
	}
	return res
}

// ToF32 casts a slice of float64 to float32. Possible loss of resolution.
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
	if n == 0 {
		return nil
	}
	d := MinD(pts...)
	res := make([]float64, d)

	// Sum
	for _, pt := range pts {
		for i, v := range pt {
			if i > d-1 {
				break
			}
			res[i] += v
		}
	}
	// Scale
	for i := 0; i < d; i++ {
		res[i] /= float64(n)
	}
	return res
}

// CrossProduct returns the cross product of the three points.
func CrossProduct(p1, p2, p3 []float64) float64 {
	return (p3[0]-p1[0])*(p2[1]-p1[1]) - (p3[1]-p1[1])*(p2[0]-p1[0])
}

// DotProduct returns the dot product of the two lines, p1-p2 and p3-p4.
func DotProduct(p1, p2, p3, p4 []float64) float64 {
	return (p2[0]-p1[0])*(p4[0]-p3[0]) + (p2[1]-p1[1])*(p4[1]-p3[1])
}

// LineAngle returns the angle of a line.
func LineAngle(p1, p2 []float64) float64 {
	return math.Atan2(p2[1]-p1[1], p2[0]-p1[0])
}

// AngleBetweenLines using Atan2 vs calculating the dot product (2xSqrt+Acos).
// Retains the directionality of the rotation from l1 to l2, unlike dot product.
func AngleBetweenLines(p1, p2, p3, p4 []float64) float64 {
	a1 := LineAngle(p1, p2)
	a2 := LineAngle(p3, p4)
	da := a2 - a1
	if da < -math.Pi {
		da += 2 * math.Pi
	} else if da > math.Pi {
		da -= 2 * math.Pi
	}
	return da
}

// MinD calculates the minimum dimensionality of point set
func MinD(pts ...[]float64) int {
	d := len(pts[0])
	for i := 1; i < len(pts); i++ {
		n := len(pts[i])
		if d > n {
			d = n
		}
	}
	return d
}
