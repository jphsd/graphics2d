package util

// TriArea returns the signed area of a triangle by finding the determinant of
// M = {{p1[0] p1[1] 1}
//      {p2[0] p2[1] 1}
//      {p3[0] p3[1] 1}}
func TriArea(p1, p2, p3 []float64) float64 {
	// Expand on third column
	det := (p2[0]*p3[1] - p3[0]*p2[1]) -
		(p1[0]*p3[1] - p3[0]*p1[1]) +
		(p1[0]*p2[1] - p2[0]*p1[1])
	return det / 2
}

// Collinear returns true if three points are on a lines (i.e. if the area of the resultant triangle is 0)
func Collinear(p1, p2, p3 []float64) bool {
	a := TriArea(p1, p2, p3)
	return Equals(0, a)
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
