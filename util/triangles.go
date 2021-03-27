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
