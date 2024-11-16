package util

import "math"

// TriArea returns the signed area of a triangle by finding the determinant of
// M = {{p1[0] p1[1] 1}
//
//	{p2[0] p2[1] 1}
//	{p3[0] p3[1] 1}}
func TriArea(p1, p2, p3 []float64) float64 {
	// Expand on third column
	det := (p2[0]*p3[1] - p3[0]*p2[1]) -
		(p1[0]*p3[1] - p3[0]*p1[1]) +
		(p1[0]*p2[1] - p2[0]*p1[1])
	return det / 2
}

var (
	SideOfLine = TriArea
)

// Left returns true if p is to the left of the line
func Left(lp1, lp2, p []float64) bool {
	return TriArea(lp1, lp2, p) > 0
}

// Left returns true if p is to the left of or on the line
func LeftOn(lp1, lp2, p []float64) bool {
	return TriArea(lp1, lp2, p) >= 0
}

// Right returns true if p is to the right of the line
func Right(lp1, lp2, p []float64) bool {
	return TriArea(lp1, lp2, p) < 0
}

// RightOn returns true if p is to the right of or on the line
func RightOn(lp1, lp2, p []float64) bool {
	return TriArea(lp1, lp2, p) <= 0
}

// Collinear returns true if three points are on a line (i.e. if the area of the resultant triangle is 0)
func Collinear(lp1, lp2, p []float64) bool {
	return Equals(TriArea(lp1, lp2, p), 0)
}

// PointInTriangle returns true if p is in the triangle formed by tp1, tp2 and tp3.
func PointInTriangle(p, tp1, tp2, tp3 []float64) bool {
	w1, w2, w3 := Barycentric(p, tp1, tp2, tp3)
	return !(w1 < 0 || w1 > 1 || w2 < 0 || w2 > 1 || w3 < 0 || w3 > 1)
}

// Barycentric converts a point on the plane into Barycentric weights given three non-coincident points, tp1, tp2 and tp3.
func Barycentric(p, tp1, tp2, tp3 []float64) (float64, float64, float64) {
	// Barycentric form (see https://en.wikipedia.org/wiki/Barycentric_coordinate_system)
	// p = d1 * t1 + d2 * t2 + d3 * t3. If p is within tri then d1 + d2 + d3 = 1
	d := tp1[0]*(tp2[1]-tp3[1]) + tp1[1]*(tp3[0]-tp2[0]) + tp2[0]*tp3[1] - tp2[1]*tp3[0]
	d = 1 / d
	d1 := (p[0]*(tp3[1]-tp1[1]) + p[1]*(tp1[0]-tp3[0]) - tp1[0]*tp3[1] + tp1[1]*tp3[0]) * d
	d2 := (p[0]*(tp2[1]-tp1[1]) + p[1]*(tp1[0]-tp2[0]) - tp1[0]*tp2[1] + tp1[1]*tp2[0]) * -d
	d3 := 1.0 - d1 - d2

	return d1, d2, d3
}

// InverseBarycentric converts Barycentric weights plus the three non-coincident points back into a point on the plane.
func InverseBarycentric(w1, w2, w3 float64, tp1, tp2, tp3 []float64) []float64 {
	return []float64{w1*tp1[0] + w2*tp2[0] + w3*tp3[0], w1*tp1[1] + w2*tp2[1] + w3*tp3[1]}
}

// Dejitter takes a list of points, greater than three long, and removes the one that forms the smallest area triangle
// with its adjacent points. If closed is true then the end points are candidates for elimination too.
func Dejitter(closed bool, pts ...[]float64) [][]float64 {
	n := len(pts)
	if n < 4 {
		return pts
	}

	min := math.MaxFloat64
	mi := 0
	if closed {
		v := TriArea(pts[n-1], pts[0], pts[1])
		if v < 0 {
			v = -v
		}
		if v < min {
			min = v
		}
		v = TriArea(pts[n-2], pts[n-1], pts[0])
		if v < 0 {
			v = -v
		}
		if v < min {
			min = v
			mi = n - 1
		}
	}
	for i := 1; i < n-1; i++ {
		v := TriArea(pts[i-1], pts[i], pts[i+1])
		if v < 0 {
			v = -v
		}
		if v < min {
			min = v
			mi = i
		}
	}
	// Remove mi'th point
	if mi == 0 {
		pts[0] = nil
		return pts[1:]
	} else if mi != n-1 {
		copy(pts[mi:], pts[mi+1:])
	}
	n--
	pts[n] = nil
	return pts[:n]
}
