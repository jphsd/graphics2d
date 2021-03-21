package util

import (
	"math"
)

// Given an arc angle (less than or equal to PI), calculate the
// points for a Bezier cubic to describe it on a circle centered
// on (0,0) with radius 1. Mid-point of the curve is (1,0)
// Error increases for values > PI/2
func CalcPointsForArc(theta float64) [][]float64 {
	phi := theta / 2
	x0 := math.Cos(phi)
	y0 := math.Sin(phi)
	x3 := x0
	y3 := -y0
	x1 := (4 - x0) / 3
	y1 := (1 - x0) * (3 - x0) / (3 * y0)
	x2 := x1
	y2 := -y1
	return [][]float64{{x3, y3}, {x2, y2}, {x1, y1}, {x0, y0}}
}

// Conversion methods for cubic Bezier to CatmullRom and v.v.
// From https://pomax.github.io/bezierinfo/#catmullconv
// p1, c1, c2, p2 => t1, p1, p2, t2
func Bezier3ToCatmul(p1, p2, p3, p4 []float64) []float64 {
	dx12 := 6 * (p1[0] - p2[0])
	dy12 := 6 * (p1[1] - p2[1])
	dx43 := 6 * (p4[0] - p3[0])
	dy43 := 6 * (p4[1] - p3[1])
	return []float64{p4[0] + dx12, p4[1] + dy12, p1[0], p1[1], p4[0], p4[1], p1[0] + dx43, p1[1] + dy43}
}

// t1, p1, p2, t2 => p1, c1, c2, p2
func CatmulToBezier3(tau float64, p1, p2, p3, p4 []float64) []float64 {
	tau *= 6
	dx31 := (p3[0] - p1[0]) / tau
	dy31 := (p3[1] - p1[1]) / tau
	dx42 := (p4[0] - p2[0]) / tau
	dy42 := (p4[1] - p2[1]) / tau
	return []float64{p2[0], p2[1], p2[0] + dx31, p2[1] + dy31, p3[0] - dx42, p3[1] - dy42, p3[0], p3[1]}
}

// B1 (flat curve) {p1, p2}
func Bezier1(pts [][]float64, t float64) []float64 {
	omt := 1 - t
	return []float64{
		omt*pts[0][0] + t*pts[1][0],
		omt*pts[0][1] + t*pts[1][1]}
}

// B2 (quad curve) {p1, c1, p2}
func Bezier2(pts [][]float64, t float64) []float64 {
	t2 := t * t
	omt := 1 - t
	omt2 := omt * omt
	omt2t := omt * 2 * t
	return []float64{
		omt2*pts[0][0] + omt2t*pts[1][0] + t2*pts[2][0],
		omt2*pts[0][1] + omt2t*pts[1][1] + t2*pts[2][1]}
}

// B3 (cubic curve) {p1, c1, c2, p2}
func Bezier3(pts [][]float64, t float64) []float64 {
	t2 := t * t
	t3 := t2 * t
	omt := 1 - t
	omt2 := omt * omt
	omt3 := omt2 * omt

	bc1 := 3 * omt2 * t
	bc2 := 3 * omt * t2
	return []float64{
		omt3*pts[0][0] + bc1*pts[1][0] + bc2*pts[2][0] + t3*pts[3][0],
		omt3*pts[0][1] + bc1*pts[1][1] + bc2*pts[2][1] + t3*pts[3][1]}
}

// DeCasteljau uses de Casteljau's algorithm for degree n curves and
// returns the point and the tangent of the line it's traversing.
// {p1, c1, c2, c3, ..., p2}
func DeCasteljau(pts [][]float64, t float64) []float64 {
	if len(pts) == 1 {
		return pts[0]
	}
	npts := make([][]float64, len(pts)-1)
	omt := 1 - t
	for i := 0; i < len(npts); i++ {
		npts[i] = []float64{
			omt*pts[i][0] + t*pts[i+1][0], omt*pts[i][1] + t*pts[i+1][1],
			pts[i+1][0] - pts[i][0], pts[i+1][1] - pts[i][1]}
	}
	return DeCasteljau(npts, t)
}

// Split curve at t into two new curves such that the end of the lhs is the
// start of the rhs
// {p1, c1, c2, c3, ..., p2}
func SplitCurve(pts [][]float64, t float64) [][][]float64 {
	n := len(pts)
	left := make([][]float64, n)
	right := make([][]float64, n)
	splitCurve(pts, n-1, 0, left, right, t)
	return [][][]float64{left, right}
}

// Helper function - note flipping of rhs
func splitCurve(pts [][]float64, nn, n int, left, right [][]float64, t float64) {
	np := len(pts)
	if np == 1 {
		left[n] = pts[0]
		right[nn-n] = pts[0]
	} else {
		np -= 1
		npts := make([][]float64, np)
		omt := 1 - t
		for i := 0; i < np; i++ {
			if i == 0 {
				left[n] = pts[0]
			}
			if i == np-1 {
				right[nn-n] = pts[np]
			}
			npts[i] = []float64{
				omt*pts[i][0] + t*pts[i+1][0], omt*pts[i][1] + t*pts[i+1][1],
				pts[i+1][0] - pts[i][0], pts[i+1][1] - pts[i][1]}
		}
		splitCurve(npts, nn, n+1, left, right, t)
	}
}

// Bezier gradient (differentiation): Order of curve drops by one and new weights
// are the difference of the original weights scaled by the original order
func CalcDerivativeWeights(w [][]float64) [][]float64 {
	n := len(w)
	res := make([][]float64, n-1)
	for i := 0; i < n-1; i++ {
		res[i] = []float64{
			float64(n) * (w[i+1][0] - w[i][0]),
			float64(n) * (w[i+1][1] - w[i][1])}
	}
	return res
}

// Bezier curve order promotion e.g. b2 to b3. Note, there's no inverse.
func CalcNextOrderWeights(w [][]float64) [][]float64 {
	n := len(w)
	k := n + 1
	ki := 1 / float64(k)
	res := make([][]float64, k)
	// First weight doesn't change
	res[0] = []float64{
		w[0][0],
		w[0][1]}
	for i := 1; i < k; i++ {
		res[i] = []float64{
			ki * (float64(k-i)*w[i][0] + float64(i)*w[i-1][0]),
			ki * (float64(k-i)*w[i][1] + float64(i)*w[i-1][1])}
	}
	return res
}

// Calculate curvature - note curve must have 2nd order derivative
// Radius of curvature at t is 1/kappa(t)
func Kappa1(dw, d2w [][]float64, t float64) float64 {
	dpt := DeCasteljau(dw, t)
	d2pt := DeCasteljau(d2w, t)
	return Kappa(dpt, d2pt)
}

// Kappa/curvature from first and second derivatives at a point
func Kappa(dpt, d2pt []float64) float64 {
	return (dpt[0]*d2pt[1] - d2pt[0]*dpt[1]) / math.Pow(dpt[0]*dpt[0]+dpt[1]*dpt[1], 1.5)
}

// Estimate kappa from three points by calculating the center of the circumcircle
func KappaC(p1, p2, p3 []float64) float64 {
	d1 := []float64{p2[0] - p1[0], p2[1] - p1[1]}
	s1 := []float64{p1[0] + d1[0]/2, p1[1] + d1[1]/2} // mid point p1-p2
	d2 := []float64{p3[0] - p2[0], p3[1] - p2[1]}
	s2 := []float64{p2[0] + d2[0]/2, p2[1] + d2[1]/2} // mid point p2-p3
	n1 := []float64{-d1[1], d1[0]}                    // normal at s1
	n2 := []float64{-d2[1], d2[0]}                    // normal at s2
	// intersection of n1 and n2
	err, ts := IntersectionTVals(s1[0], s1[1], s1[0]+n1[0], s1[1]+n1[1], s2[0], s2[1], s2[0]+n2[0], s2[1]+n2[1])
	if err != nil {
		// p1, p2 and p3 are coincident
		return 0
	}
	c := []float64{s1[0] + ts[0]*n1[0], s1[1] + ts[0]*n1[1]} // cc center
	dx := p1[0] - c[0]
	dy := p1[1] - c[1]
	r := math.Sqrt(dx*dx + dy*dy) // distance p1-c
	if ts[0] < 0 {
		return -1 / r
	}
	return 1 / r
}

// Menger curvature: 4*area / (d(p1, p2).d(p2s, p3).d(p3, p1))
// Same result as above but with more square roots...
func KappaM(p1, p2, p3 []float64) float64 {
	a := TriArea(p1, p2, p3)
	denom := DistanceE(p1, p2) * DistanceE(p2, p3) * DistanceE(p3, p1)
	if Equals(denom, 0) {
		// p1, p2 and p3 are coincident
		return 0
	}
	return 4 * a / denom
}
