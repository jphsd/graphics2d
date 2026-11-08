package graphics2d

import (
	"fmt"
	"math"
	"strings"

	"github.com/jphsd/graphics2d/util"
)

// Part represents a path step with the previous step's end point prepended.
type Part [][]float64
type Parts []Part

// String converts a part into a string.
func (p Part) String() string {
	var res strings.Builder
	for i, pt := range p {
		if i != 0 {
			res.WriteString(" ")
		}
		res.WriteString(fmt.Sprintf("%.2f,%.2f", pt[0], pt[1]))
	}
	return res.String()
}

// PartsToPath constructs a new path by concatenating the parts.
func PartsToPath(parts ...Part) *Path {
	if len(parts) == 0 {
		return nil
	}

	res := NewPath(parts[0][0])
	if len(parts[0]) == 1 {
		return res
	}

	for i, part := range parts {
		for j, pt := range part {
			if InvalidPoint(pt) {
				panic(fmt.Sprintf("invalid point (NaN) in part %d,%d", i, j))
			}
		}
		res.AddStep(part[1:]...)
	}
	return res
}

// flattenPart successively splits the part until the control points are within d2 of the line
// joing the part start with the part end. Returns a list of line parts.
func (part Part) flatten(d2 float64) Parts {
	if part.cpWithinD2(d2) {
		return Parts{{part[0], part[len(part)-1]}}
	}
	lr := util.SplitCurve(part, 0.5)
	res := append(Parts{}, Part(lr[0]).flatten(d2)...)
	res = append(res, Part(lr[1]).flatten(d2)...)
	return res
}

func (part Part) cpWithinD2(d2 float64) bool {
	// First and last are end points
	l := len(part)
	if l == 2 {
		// Trivial case
		return true
	}
	start, cpts, end := part[0], part[1:l-1], part[l-1]
	for _, cp := range cpts {
		pd2, _, _ := util.DistanceToLineSquared(start, end, cp)
		if pd2 > d2 {
			return false
		}
	}

	return true
}

// FlattenPart works by subdividing the curve until its control points are within d^2 (d squared)
// of the line through the end points.
func (part Part) Flatten(d float64) Parts {
	return part.flatten(d * d)
}

// PartLength returns the approximate length of a part by flattening it to the supplied degree of flatness.
func (part Part) Length(d float64) float64 {
	parts := part.flatten(d * d)
	sum := 0.0
	for _, part := range parts {
		sum += util.DistanceE(part[0], part[1])
	}
	return sum
}

// SimplifyExtremities chops curve into pieces based on maxima, minima and inflections in x and y.
func (part Part) SimplifyExtremities() Parts {
	tvals := util.CalcExtremities(part)
	nt := len(tvals)
	if nt < 3 {
		return Parts{part}
	}
	// Convert tvals to relative tvals
	rtvals := make([]float64, nt)
	lt := 0.0
	ll := 1.0
	for i := range nt {
		if i == 0 {
			// Reset
			lt = tvals[i]
			ll = 1 - lt
			rtvals[i] = lt
			continue
		}
		t := tvals[i]
		d := t - lt
		rtvals[i] = d / ll
		lt = t
		ll -= d
	}

	rhs := part
	res := make(Parts, nt-1)
	for i := 1; i < nt-1; i++ {
		lr := util.SplitCurve(rhs, rtvals[i])
		res[i-1] = lr[0]
		rhs = lr[1]
	}
	res[nt-2] = rhs
	return res
}

// SimplifyPart recursively cuts the curve in half until CPSafe is
// satisfied.
func (part Part) Simplify() Parts {
	if part.CPSafe() {
		return Parts{part}
	}
	lr := util.SplitCurve(part, 0.5)
	res := append(Parts{}, Part(lr[0]).Simplify()...)
	res = append(res, Part(lr[1]).Simplify()...)
	return res
}

// SafeFraction if greater than 0 causes Simplify to perform a check of the mid-point against
// the part centroid. If the two are within SafeFraction of the distance from p[0] to the centroid
// then no further subdivision of the curve is performed.
var SafeFraction float64 = -1

// CPSafe returns true if all the control points are on the same side of
// the line formed by start and the last part points and the point at t = 0.5 is close
// to the centroid of the part.
func (part Part) CPSafe() bool {
	n := len(part)
	if n < 3 {
		// Either a point or line
		return true
	}

	start := part[0]
	end := part[n-1]
	side := util.CrossProduct(start, end, part[1]) < 0
	for i := 2; i < n-1; i++ {
		if (util.CrossProduct(start, end, part[i]) < 0) != side {
			return false
		}
	}

	if n == 3 {
		return true
	}

	if SafeFraction > 0 {
		// Check mid-point against centroid
		// scale against distance between p0 and centroid
		centroid := util.Centroid(part...)
		hp := util.DeCasteljau(part, 0.5)
		p0dx := centroid[0] - part[0][0]
		p0dy := centroid[1] - part[0][1]
		p0ds := p0dx*p0dx + p0dy*p0dy
		hpdx := centroid[0] - hp[0]
		hpdy := centroid[1] - hp[1]
		hpds := hpdx*hpdx + hpdy*hpdy
		return math.Sqrt(hpds) < math.Sqrt(p0ds)*SafeFraction
	}

	return true
}

// Reverse reverses the order (and points) of the supplied part slice.
func (parts Parts) Reverse() Parts {
	n := len(parts)
	res := make(Parts, n)
	for i, j := 0, n-1; i < n; i++ {
		res[i] = parts[j].ReversePoints()
		j--
	}
	return res
}

// [pts][x/y]
func (cp Part) ReversePoints() Part {
	n := len(cp)
	res := make(Part, n)
	for i, j := 0, n-1; i < n; i++ {
		res[i] = cp[j]
		j--
	}
	return res
}

// PartsIntersection returns the location of where the two parts intersect or nil. Assumes the parts
// are the result of simplification. Uses a brute force approach for curves with d as the flattening
// value.
func PartsIntersection(part1, part2 Part, d float64) []float64 {
	// Test bounding boxes first
	bb1, bb2 := util.BoundingBox(part1...), util.BoundingBox(part2...)
	if !util.BBOverlap(bb1, bb2) {
		return nil
	}

	// Flatten and calculate bounding boxes for part2
	fparts1, fparts2 := part1.Flatten(d), part2.Flatten(d)
	bbs2 := make([][][]float64, len(fparts2))
	for i, part := range fparts2 {
		bbs2[i] = util.BoundingBox(part...)
	}

	// Test each line in part1 against lines in part2 until we find an intersection
	for _, part := range fparts1 {
		bb1 = util.BoundingBox(part...)
		s1, e1 := part[0], part[len(part)-1]
		for j, bbp2 := range bbs2 {
			if !util.BBOverlap(bb1, bbp2) {
				continue
			}
			// Bounding boxes of lines overlap - see if they intersect
			s2, e2 := fparts2[j][0], fparts2[j][len(fparts2[j])-1]
			tvals, err := util.IntersectionTValsP(s1, e1, s2, e2)
			if err != nil || tvals[0] < 0 || tvals[0] > 1 || tvals[1] < 0 || tvals[1] > 1 {
				continue
			}
			//return []float64{util.Lerp(tvals[0], s1[0], e1[0]), util.Lerp(tvals[0], s1[1], e1[1])}
			return Lerp(tvals[0], s1, e1)
		}
	}

	return nil
}

// Returns closest point on part to pt, dist2 and t [0-1]
func (part Part) bs(pt []float64, ts, ds, te, de float64) ([]float64, float64, float64) {
	dt := te - ts
	if dt < 0.00001 {
		prj := util.DeCasteljau(part, ts)
		return prj, ds, ts
	}

	t := []float64{ts, ts + dt/4, ts + dt/2, ts + 3*dt/4, te}
	dl := dist2(pt, util.DeCasteljau(part, t[1]))
	dm := dist2(pt, util.DeCasteljau(part, t[2]))
	dr := dist2(pt, util.DeCasteljau(part, t[3]))
	d := []float64{ds, dl, dm, dr, de}
	nd := len(d)

	ci, cd := 0, d[0]
	for i := 1; i < nd; i++ {
		if d[i] < cd {
			ci = i
			cd = d[i]
		}
	}
	if ci == 0 {
		// Only search t[0] to t[1]
		return part.bs(pt, t[0], d[0], t[1], d[1])
	}
	if ci == nd-1 {
		// Only search t[3] to t[4]
		return part.bs(pt, t[3], d[3], t[4], d[4])
	}

	// Search both t[ci-1] to t[ci] and t[ci] to t[ci+1]
	pl, d1, tl := part.bs(pt, t[ci-1], d[ci-1], t[ci], d[ci])
	pr, d2, tr := part.bs(pt, t[ci], d[ci], t[ci+1], d[ci+1])
	if d1 < d2 {
		return pl, d1, tl
	}
	return pr, d2, tr
}

// Lerp performs a linear interpolation between two points.
func Lerp(t float64, p1, p2 []float64) []float64 {
	return []float64{util.Lerp(t, p1[0], p2[0]), util.Lerp(t, p1[1], p2[1])}
}
