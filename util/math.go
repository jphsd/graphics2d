package util

import (
	"fmt"
	"math"
)

const (
	Pi = math.Pi
	TwoPi = Pi * 2
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

// EqualsP returns true if two points are equal (upto the extent of shared dimensions).
func EqualsP(v1, v2 []float64) bool {
	lv1, lv2 := len(v1), len(v2)
	min := lv1
	if lv2 < min {
		min = lv2
	}
	for i := 0; i < min; i++ {
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
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	return math.Hypot(dx, dy)
}

// LineNormalToPoint returns the point on the line that's normal to the specified point, in
// absolute terms and as a t value. The distance from the point to the line is returned too.
func LineNormalToPoint(lp1, lp2, p []float64) ([]float64, float64, float64) {
	d2, pt, t := DistanceToLineSquared(lp1, lp2, p)
	return pt, t, math.Sqrt(d2)
}

// DistanceToLineSquared calculates the squared Euclidean length of the normal from a point to
// the line. Returns the distance squared, the line intercept and the t value.
func DistanceToLineSquared(lp1, lp2, p []float64) (float64, []float64, float64) {
	dx := lp2[0] - lp1[0]
	dy := lp2[1] - lp1[1]
	// Check for line degeneracy
	if Equals(0, dx) && Equals(0, dy) {
		return DistanceESquared(lp1, p), lp1, 0
	}
	qx := p[0] + dy
	qy := p[1] - dx
	ts, err := IntersectionTVals(lp1[0], lp1[1], lp2[0], lp2[1], p[0], p[1], qx, qy)
	if err != nil {
		return DistanceESquared(lp1, p), lp1, 0
	}
	t := ts[1]
	ip := []float64{Lerp(t, p[0], qx), Lerp(t, p[1], qy)}
	dx = ip[0] - p[0]
	dy = ip[1] - p[1]
	return dx*dx + dy*dy, ip, ts[0]
}

// PointOnLine returns if p is on the line and t for p if it is.
func PointOnLine(lp1, lp2, p []float64) (bool, float64) {
	cp := CrossProduct(lp1, lp2, p)
	if !Equals(cp, 0) {
		return false, -1
	}
	dx := lp2[0] - lp1[0]
	if Equals(dx, 0) {
		// Vertical
		return true, (p[1] - lp1[1]) / (lp2[1] - lp2[0])
	}
	return true, (p[0] - lp1[0]) / dx
}

// True if lines are parallel
func Parallel(lp1, lp2, lp3, lp4 []float64) bool {
	_, err := IntersectionTValsP(lp1, lp2, lp3, lp4)
	return err != nil
}

// True if lines are on the same infinite line
func Coincident(lp1, lp2, lp3, lp4 []float64) bool {
	if !Parallel(lp1, lp2, lp3, lp4) {
		return false
	}
	cp := CrossProduct(lp1, lp2, lp3)
	return Equals(cp, 0)
}

// ToF64 casts a slice of float32 to float64.
func ToF64(pts ...float32) []float64 {
	res := make([]float64, len(pts))
	for i, v := range pts {
		if v != v {
			panic("casting NaN")
		}
		res[i] = float64(v)
	}
	return res
}

// ToF32 casts a slice of float64 to float32. Possible loss of resolution.
func ToF32(pts ...float64) []float32 {
	res := make([]float32, len(pts))
	for i, v := range pts {
		if v != v {
			panic("casting NaN")
		}
		res[i] = float32(v)
	}
	return res
}

// Centroid returns the vertex centroid of a set of points.
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
// Since the inputs are all in the x-y plane (i.e. z = 0), only the magnitude of the
// resultant z vector is returned (the x and y vectors are both 0).
func CrossProduct(p1, p2, p3 []float64) float64 {
	return (p3[0]-p1[0])*(p2[1]-p1[1]) - (p3[1]-p1[1])*(p2[0]-p1[0])
}

// DotProduct returns the dot product of the two lines, p1-p2 and p3-p4.
func DotProduct(p1, p2, p3, p4 []float64) float64 {
	return (p2[0]-p1[0])*(p4[0]-p3[0]) + (p2[1]-p1[1])*(p4[1]-p3[1])
}

// Vec returns the vector joining two points.
func Vec(p1, p2 []float64) []float64 {
	return []float64{p2[0] - p1[0], p2[1] - p1[1]}
}

// VecMag returns the magnitude of the vector.
func VecMag(v []float64) float64 {
	return math.Hypot(v[0], v[1])
}

// VecNormalize scales a vector to unit length.
func VecNormalize(v []float64) []float64 {
	d := VecMag(v)
	return []float64{v[0] / d, v[1] / d}
}

// LineAngle returns the angle of a line. Result is in range [-Pi, Pi]
func LineAngle(p1, p2 []float64) float64 {
	return math.Atan2(p2[1]-p1[1], p2[0]-p1[0])
}

// AngleBetweenLines using Atan2 vs calculating the dot product (2xSqrt+Acos).
// Retains the directionality of the rotation from l1 to l2, unlike dot product.
// The result is in the range [-Pi,Pi].
func AngleBetweenLines(p1, p2, p3, p4 []float64) float64 {
	a1 := LineAngle(p1, p2)
	a2 := LineAngle(p3, p4)
	da := a2 - a1
	if da < -Pi {
		da += TwoPi
	} else if da > Pi {
		da -= TwoPi
	}
	return da
}

// DotProductAngle returns the angle between two lines using the dot product method. The
// result is in the range [0,Pi].
func DotProductAngle(p1, p2, p3, p4 []float64) float64 {
	v1 := VecNormalize(Vec(p1, p2))
	v2 := VecNormalize(Vec(p3, p4))
	dp := v1[0]*v2[0] + v1[1]*v2[1]
	// Clamp rounding errors
	if dp < -1 {
		dp = -1
	} else if dp > 1 {
		dp = 1
	}
	return math.Acos(dp)
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

// Circumcircle returns the circle (center and radius) that passes through the three points.
func Circumcircle(p1, p2, p3 []float64) []float64 {
	// Translate p1, p2 and p3 s.t. p1 is at the origin
	b := []float64{p2[0] - p1[0], p2[1] - p1[1]}
	b2 := b[0]*b[0] + b[1]*b[1]
	c := []float64{p3[0] - p1[0], p3[1] - p1[1]}
	c2 := c[0]*c[0] + c[1]*c[1]
	d := 2 * (b[0]*c[1] - c[0]*b[1])
	d = 1 / d
	u := []float64{(c[1]*b2 - b[1]*c2) * d, (b[0]*c2 - c[0]*b2) * d}
	r2 := u[0]*u[0] + u[1]*u[1]
	return []float64{u[0] + p1[0], u[1] + p1[1], math.Sqrt(r2)}
}

// AngleInRange returns true if angle b is within [a,a+r].
func AngleInRange(a, r, b float64) bool {
	// Ensure a, b in [-pi,pi], r in [-2pi,2pi]
	a, b = mapAngle(a), mapAngle(b)
	for r < -TwoPi {
		r += TwoPi
	}
	for r > TwoPi {
		r -= TwoPi
	}

	ob := b - a
	if r < 0 {
		return !(ob < r || ob > 0)
	}
	return !(ob < 0 || ob > r)
}

func mapAngle(a float64) float64 {
	for a < -Pi {
		a += TwoPi
	}
	for a > Pi {
		a -= TwoPi
	}
	return a
}

// Lerp returns the value (1-t)*start + t*end.
func Lerp(t, start, end float64) float64 {
	return (1-t)*start + t*end
}
