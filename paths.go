package graphics2d

import (
	"math"
	"math/rand"

	"github.com/jphsd/graphics2d/util"
)

// Mathematical constant.
const (
	TwoPi = 2 * math.Pi
)

// A collection of part and path creation functions.

// MakeArcParts creates at least one cubic bezier that describes a curve from offs to
// offs+ang centered on {cx, cy} with radius r.
func MakeArcParts(cx, cy, r, offs, ang float64) [][][]float64 {
	if util.Equals(ang, 0) {
		// Return just a point
		pt := []float64{cx + r*math.Cos(offs), cy + r*math.Sin(offs)}
		return [][][]float64{{pt, pt}}
	}

	a := ang
	rev := ang < 0
	if rev {
		a = -a
	}

	// Calculate number of curves to create - necessary since curve errors
	// are apparent for angles > Pi/2
	n := 1
	for true {
		if a < math.Pi/2 {
			break
		}
		a /= 2
		n *= 2
	}
	cp := util.CalcPointsForArc(a)

	if rev {
		cp = [][]float64{cp[3], cp[2], cp[1], cp[0]}
		a = -a
	}

	res := make([][][]float64, n)
	for i := 0; i < n; i++ {
		xfm := CreateTransform(cx, cy, r, offs+a/2)
		res[i] = xfm.Apply(cp...)
		offs += a
	}
	return res
}

// Point returns a path containing the point.
func Point(pt []float64) *Path {
	return NewPath(pt)
}

// Line returns a path describing the line.
func Line(pt1, pt2 []float64) *Path {
	np := NewPath(pt1)
	np.AddStep(pt2)
	return np
}

// PolyLine returns a path with lines joining successive points.
func PolyLine(pts ...[]float64) *Path {
	if len(pts) == 0 {
		return nil
	}
	np := NewPath(pts[0])
	for i := 1; i < len(pts); i++ {
		np.AddStep(pts[i])
	}
	return np
}

// Polygon returns a closed path with lines joining successive points.
func Polygon(pts ...[]float64) *Path {
	if len(pts) == 0 {
		return nil
	}
	np := NewPath(pts[0])
	for i := 1; i < len(pts); i++ {
		np.AddStep(pts[i])
	}
	np.Close()
	return np
}

// Curve returns a path describing the polynomial curve.
func Curve(pts ...[]float64) *Path {
	if len(pts) == 0 {
		return nil
	}
	np := NewPath(pts[0])
	np.AddStep(pts[1:]...)
	return np
}

// PolyCurve returns a path describing the polynomial curves.
func PolyCurve(pts ...[][]float64) *Path {
	if len(pts) == 0 {
		return nil
	}
	np := NewPath(pts[0][0])
	for i := 0; i < len(pts); i++ {
		np.AddStep(pts[i][1:]...)
	}
	return np
}

// ArcStyle defines the type of arc - open, chord (closed) and pie (closed).
type ArcStyle int

// Constants for arc styles.
const (
	ArcOpen ArcStyle = iota
	ArcChord
	ArcPie
)

// Arc returns a path with an arc centered on c with radius r from offs in the direction and length of ang.
func Arc(c []float64, r, offs, ang float64, s ArcStyle) *Path {
	// Limit offs and ang to +/- 2 pi
	for offs > TwoPi {
		offs -= TwoPi
	}
	for offs < -TwoPi {
		offs += TwoPi
	}
	for ang > TwoPi {
		ang -= TwoPi
	}
	for ang < -TwoPi {
		ang += TwoPi
	}

	parts := MakeArcParts(c[0], c[1], r, offs, ang)
	np := PartsToPath(parts...)
	switch s {
	case ArcChord:
		np.Close()
		return np
	case ArcPie:
		np.AddStep(c)
		np.Close()
		return np
	}
	return np
}

// ArcFromPoint returns a path describing an arc starting from pt based on c and ang.
func ArcFromPoint(pt, c []float64, ang float64, s ArcStyle) *Path {
	dx := pt[0] - c[0]
	dy := pt[1] - c[1]
	r := math.Hypot(dx, dy)
	offs := math.Atan2(dy, dx)
	return Arc(c, r, offs, ang, s)
}

// ArcFromPoints returns a path describing an arc passing through a, b and c such that the
// arc starts at a, passes through b and ends at c.
func ArcFromPoints(a, b, c []float64, s ArcStyle) *Path {
	cp := util.Circumcircle(a, b, c)
	if math.IsInf(cp[2], 0) {
		return Line(a, c)
	}
	aa, ca := util.LineAngle(cp, a), util.LineAngle(cp, c)
	ang := ca - aa // [-2pi,2pi]

	var oang float64 // [-2pi,2pi] - the opposite of ang such that |ang| + |oang| = 2pi
	if ang < 0 {
		oang = TwoPi + ang
	} else {
		oang = ang - TwoPi
	}
	agto := math.Abs(ang) > math.Abs(oang)

	// if cp in tri, |ang| > pi, else |ang| < pi
	if util.PointInTriangle(cp, a, b, c) {
		// Choose larger of |ang|, |oang|
		if agto {
			return Arc(cp, cp[2], aa, ang, s)
		} else {
			return Arc(cp, cp[2], aa, oang, s)
		}
	} else {
		// Choose smaller of |ang|, |oang|
		if agto {
			return Arc(cp, cp[2], aa, oang, s)
		} else {
			return Arc(cp, cp[2], aa, ang, s)
		}
	}
}

// PolyArcFromPoint returns a path concatenating the arcs.
func PolyArcFromPoint(pt []float64, cs [][]float64, angs []float64) *Path {
	n, na := len(cs), len(angs)
	if na < n {
		n = na
	}

	parts := [][][]float64{}
	cp := pt
	for i := 0; i < n; i++ {
		dx := cp[0] - cs[i][0]
		dy := cp[1] - cs[i][1]
		r := math.Hypot(dx, dy)
		offs := math.Atan2(dy, dx)
		tmp := MakeArcParts(cs[i][0], cs[i][1], r, offs, angs[i])
		last := tmp[len(tmp)-1]
		cp = last[len(last)-1]
		parts = append(parts, tmp...)
	}

	res := PartsToPath(parts...)
	return res
}

// Circle returns a closed path describing a circle centered on c with radius r.
func Circle(c []float64, r float64) *Path {
	ax, ay := c[0], c[1]
	np := PartsToPath(MakeArcParts(ax, ay, r, 0, TwoPi)...)
	np.Close()
	return np
}

// RegularPolygon returns a closed path describing an n-sided polygon given the initial edge.
func RegularPolygon(pt1, pt2 []float64, n int) *Path {
	da := TwoPi / float64(n)
	cosDa, sinDa := math.Cos(da), math.Sin(da)
	dx, dy := pt2[0]-pt1[0], pt2[1]-pt1[1]
	np := NewPath(pt1)
	cp := pt2
	np.AddStep(cp)
	for i := 1; i < n-1; i++ {
		ncp := []float64{cp[0] + dx*cosDa - dy*sinDa, cp[1] + dx*sinDa + dy*cosDa}
		np.AddStep(ncp)
		dx, dy = ncp[0]-cp[0], ncp[1]-cp[1]
		cp = ncp
	}
	np.Close()
	return np
}

// ReentrantPolygon returns a closed path describing an n pointed star.
func ReentrantPolygon(c []float64, r float64, n int, t, ang float64) *Path {
	ang -= math.Pi / 2 // So ang = 0 has the start of the polygon pointing up
	da := TwoPi / float64(n)
	cosDa, sinDa := math.Cos(da), math.Sin(da)
	ri := r * math.Cos(da/2) * t
	skip := util.Equals(t, 1)
	dxe, dye := r*math.Cos(ang), r*math.Sin(ang)
	dxi, dyi := ri*math.Cos(ang+da/2), ri*math.Sin(ang+da/2)
	np := NewPath([]float64{c[0] + dxe, c[1] + dye})
	dxe, dye = dxe*cosDa-dye*sinDa, dxe*sinDa+dye*cosDa
	for i := 0; i < n; i++ {
		if !skip {
			np.AddStep([]float64{c[0] + dxi, c[1] + dyi})
			dxi, dyi = dxi*cosDa-dyi*sinDa, dxi*sinDa+dyi*cosDa
		}
		np.AddStep([]float64{c[0] + dxe, c[1] + dye})
		dxe, dye = dxe*cosDa-dye*sinDa, dxe*sinDa+dye*cosDa
	}
	np.Close()
	return np
}

// IrregularPolygon returns an n sided polgon guaranteed to be located within a circle of radius r centered on cp.
// If nr is set to true then polygon is forced to be non-reentrant.
func IrregularPolygon(cp []float64, r float64, n int, nr bool) *Path {
	if n < 3 {
		n = 3
	}
	tinc := TwoPi / float64(n)
	toffs := TwoPi * rand.Float64()
	fs := make([]float64, n)
	rs := make([][]float64, n)
	ps := make([][]float64, n)
	for i := 0; i < n; i++ {
		f := rand.Float64()
		if f < 0.1 {
			f = 0.1
		}
		fs[i] = f
		xr, yr := math.Cos(toffs)*r, math.Sin(toffs)*r
		rs[i] = []float64{xr + cp[0], yr + cp[1]}
		ps[i] = []float64{xr*f + cp[0], yr*f + cp[1]}
		toffs += tinc
	}

	// Iterate until none are reentrant
	nrc := 0
	cur := 0
	for nr && nrc < n {
		// See where intersection lies
		var pre, post []float64
		if cur == 0 {
			pre = ps[n-1]
		} else {
			pre = ps[cur-1]
		}
		if cur == n-1 {
			post = ps[0]
		} else {
			post = ps[cur+1]
		}
		isect, _ := util.IntersectionTValsP(cp, rs[cur], pre, post)
		if isect[0] > fs[cur] {
			// Move point outwards
			fs[cur] = isect[0] + 0.1
			if fs[cur] > 1 {
				fs[cur] = 1
			}
			ps[cur] = []float64{(rs[cur][0]-cp[0])*fs[cur] + cp[0], (rs[cur][1]-cp[1])*fs[cur] + cp[1]}
			cur++
			if cur == n {
				cur = 0
			}
			nrc = 0
			continue
		}
		cur++
		if cur == n {
			cur = 0
		}
		nrc++
	}
	path := NewPath(ps[0])
	for i := 1; i < n; i++ {
		path.AddStep(ps[i])
	}
	path.Close()
	return path
}

// Lune returns a closed path made up of two arcs with end points at c plus/minus r0 in y, all rotated by th.
// The arcs are calculated from the circumcircles of the two triangles defined by the end points, and c displaced
// by r1 or r2 in x.
func Lune(c []float64, r0, r1, r2, th float64) *Path {
	a, b := []float64{c[0], c[1] + r0}, []float64{c[0], c[1] - r0}

	var p1, p2 *Path

	if util.Equals(r1, 0) {
		p1 = Line(a, b)
	} else {
		d := []float64{c[0] + r1, c[1]}
		cc := util.Circumcircle(a, b, d)
		var dx float64
		if c[0] < cc[0] {
			dx = cc[0] - c[0]
		} else {
			dx = c[0] - cc[0]
		}
		ang := 2 * math.Atan(r0/dx)
		// min or maj ang - Does cc lie between c and d?
		if (c[0] < d[0] && c[0] < cc[0]) || (c[0] > d[0] && cc[0] < c[0]) {
			ang = TwoPi - ang
		}
		// Ang direction
		if cc[0] < d[0] {
			ang = -ang
		}
		p1 = ArcFromPoint(a, cc, ang, ArcOpen)
	}

	if util.Equals(r2, 0) {
		p2 = Line(b, a)
	} else {
		d := []float64{c[0] + r2, c[1]}
		cc := util.Circumcircle(a, b, d)
		var dx float64
		if c[0] < cc[0] {
			dx = cc[0] - c[0]
		} else {
			dx = c[0] - cc[0]
		}
		ang := 2 * math.Atan(r0/dx)
		// min or maj ang - Does cc lie between c and d?
		if (c[0] < d[0] && c[0] < cc[0]) || (c[0] > d[0] && cc[0] < c[0]) {
			ang = TwoPi - ang
		}
		// Ang direction
		if cc[0] > d[0] {
			ang = -ang
		}
		p2 = ArcFromPoint(b, cc, ang, ArcOpen)
	}

	p1.Concatenate(p2)
	p1.Close()
	xfm := NewAff3()
	xfm.RotateAbout(th, c[0], c[1])
	p1 = p1.Transform(xfm)
	return p1
}

// Square returns a closed path describing a square with side s, centered on c.
func Square(c []float64, s float64) *Path {
	hs := s / 2
	sx, sy := c[0]-hs, c[1]-hs
	points := [][]float64{
		{sx, sy},
		{sx + s, sy},
		{sx + s, sy + s},
		{sx, sy + s},
	}
	return Polygon(points...)
}

// Rectangle returns a closed path describing a rectangle with sides w and h, centered on c.
func Rectangle(c []float64, w, h float64) *Path {
	hw, hh := w/2, h/2
	sx, sy := c[0]-hw, c[1]-hh
	points := [][]float64{
		{sx, sy},
		{sx + w, sy},
		{sx + w, sy + h},
		{sx, sy + h},
	}
	return Polygon(points...)
}

// Sqrt3 is the square root of 3
const Sqrt3 = 1.7320508075688772935274463415058723669428052538103806280558069794519330169088

// Equilateral returns a closed path describing an equliateral triangle with side s, centered on c.
func Equilateral(c []float64, s float64) *Path {
	sx, sy := c[0], c[1]-s/Sqrt3
	hs := s / 2
	dy := hs * Sqrt3
	points := [][]float64{
		{sx, sy},
		{sx + hs, sy + dy},
		{sx - hs, sy + dy},
	}
	return Polygon(points...)
}
