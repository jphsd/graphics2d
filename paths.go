package graphics2d

import (
	"math"
	"math/rand"

	"github.com/jphsd/graphics2d/util"
)

// Mathematical constants.
const (
	Pi     = math.Pi
	TwoPi  = 2 * Pi
	HalfPi = Pi / 2
	Sqrt3  = 1.7320508075688772935274463415058723669428052538103806280558069794519330169088
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
		if a < HalfPi {
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
		xfm := CreateAffineTransform(cx, cy, r, offs+a/2)
		res[i] = xfm.Apply(cp...)
		offs += a
	}
	return res
}

// MakeRoundedParts uses the tangents p1-p2 and p2-p3, and the radius r to figure an arc between them.
func MakeRoundedParts(p1, p2, p3 []float64, r float64) [][][]float64 {
	theta, _, _ := util.AngleBetweenLines(p1, p2, p3, p2)
	neg := theta < 0
	if neg {
		theta = -theta
	}
	t2 := theta / 2
	tt2 := math.Tan(t2)

	// Check r is < min(p12, p23) / 2
	v1, v2 := util.Vec(p1, p2), util.Vec(p2, p3)
	d1, d2 := util.VecMag(v1), util.VecMag(v2)
	m1, m2 := d1, d2
	md := min(m1, m2)
	r = min(tt2*md, r)

	// Find intersection of arc with p1-p2
	u1 := []float64{v1[0] / d1, v1[1] / d1}
	s := r / tt2
	i12 := []float64{p2[0] - s*u1[0], p2[1] - s*u1[1]}

	// Calc center
	c := []float64{i12[0], i12[1]}
	theta = Pi - theta
	n1 := []float64{u1[1], -u1[0]}
	if neg {
		n1[0], n1[1] = -n1[0], -n1[1]
	} else {
		theta = -theta
	}
	c = []float64{c[0] + r*n1[0], c[1] + r*n1[1]}

	// Calc offset
	a12 := math.Atan2(-n1[1], -n1[0])
	return MakeArcParts(c[0], c[1], r, a12, theta)
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
	for i := range len(pts) {
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
	aa, ba, ca := util.LineAngle(cp, a), util.LineAngle(cp, b), util.LineAngle(cp, c)
	ang := ca - aa // [-2pi,2pi]

	if util.AngleInRange(aa, ang, ba) {
		return Arc(cp, cp[2], aa, ang, s)
	}

	var oang float64 // [-2pi,2pi] - the opposite of ang such that |ang| + |oang| = 2pi
	if ang < 0 {
		oang = TwoPi + ang
	} else {
		oang = ang - TwoPi
	}
	return Arc(cp, cp[2], aa, oang, s)
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

// RegularPolygon returns a closed path describing an n-sided polygon centered on c
// rotated by th. th = 0 => polygon sits on its base.
func RegularPolygon(n int, c []float64, s, th float64) *Path {
	if n < 3 {
		n = 3
	}
	a := Pi / float64(n)
	r := s / (2 * math.Sin(a))
	sa := HalfPi + a + th
	da := 2 * a
	points := make([][]float64, n)
	for i, _ := range points {
		points[i] = []float64{r*math.Cos(sa) + c[0], r*math.Sin(sa) + c[1]}
		sa += da
	}
	return Polygon(points...)
}

// ReentrantPolygon returns a closed path describing an n pointed star.
func ReentrantPolygon(c []float64, r float64, n int, t, ang float64) *Path {
	ang -= HalfPi // So ang = 0 has the start of the polygon pointing up
	da := TwoPi / float64(n)
	cosDa, sinDa := math.Cos(da), math.Sin(da)
	ri := r * math.Cos(da/2) * t
	skip := util.Equals(t, 1)
	dxe, dye := r*math.Cos(ang), r*math.Sin(ang)
	dxi, dyi := ri*math.Cos(ang+da/2), ri*math.Sin(ang+da/2)
	np := NewPath([]float64{c[0] + dxe, c[1] + dye})
	dxe, dye = dxe*cosDa-dye*sinDa, dxe*sinDa+dye*cosDa
	for range n {
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
func IrregularPolygon(c []float64, r float64, n int, nr bool) *Path {
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
		rs[i] = []float64{xr + c[0], yr + c[1]}
		ps[i] = []float64{xr*f + c[0], yr*f + c[1]}
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
		isect, _ := util.IntersectionTValsP(c, rs[cur], pre, post)
		if isect[0] > fs[cur] {
			// Move point outwards
			fs[cur] = isect[0] + 0.1
			if fs[cur] > 1 {
				fs[cur] = 1
			}
			ps[cur] = []float64{(rs[cur][0]-c[0])*fs[cur] + c[0], (rs[cur][1]-c[1])*fs[cur] + c[1]}
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
	p1 = p1.Process(&TransformProc{xfm})[0]
	return p1
}

// Lune2 returns a closed path made up of two arcs with the supplied end points.
// The arcs are calculated from the circumcircles of the two triangles defined by the end points,
// and their midpoint displaced by r1 or r2.
func Lune2(p1, p2 []float64, r1, r2 float64) *Path {
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	th := math.Atan2(dy, dx) + HalfPi
	d := math.Hypot(dx, dy)

	return Lune([]float64{p1[0] + dx/2, p1[1] + dy/2}, d/2, r1, r2, th)
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

// ExtendLine returns the line that passes through the bounds (or nil) defined by the line equation of
// pt1 and pt2.
func ExtendLine(pt1, pt2 []float64, bounds [][]float64) *Path {
	//IntersectionTValsP(p1, p2, p3, p4 []float64) ([]float64, error)
	rp1 := []float64{}
	set := false

	// Top
	b1, b2 := bounds[0], []float64{bounds[0][0], bounds[1][1]}
	tvals, err := util.IntersectionTValsP(b1, b2, pt1, pt2)
	if err != nil {
		return nil
	}
	t := tvals[0]
	if t > 0 && t < 1 {
		ont := 1 - t
		rp1 = []float64{ont*b1[0] + t*b2[0], ont*b1[1] + t*b2[1]}
		set = true
	}

	// LHS
	b2 = []float64{bounds[1][0], bounds[0][1]}
	tvals, err = util.IntersectionTValsP(b1, b2, pt1, pt2)
	if err != nil {
		return nil
	}
	t = tvals[0]
	if t > 0 && t < 1 {
		ont := 1 - t
		tmp := []float64{ont*b1[0] + t*b2[0], ont*b1[1] + t*b2[1]}
		if set {
			return Line(rp1, tmp)
		}
		rp1 = tmp
		set = true
	}

	// RHS
	b1, b2 = bounds[1], []float64{bounds[1][0], bounds[0][1]}
	tvals, err = util.IntersectionTValsP(b1, b2, pt1, pt2)
	if err != nil {
		return nil
	}
	t = tvals[0]
	if t > 0 && t < 1 {
		ont := 1 - t
		tmp := []float64{ont*b1[0] + t*b2[0], ont*b1[1] + t*b2[1]}
		if set {
			return Line(rp1, tmp)
		}
		rp1 = tmp
	}

	// Bottom
	b2 = []float64{bounds[0][0], bounds[1][1]}
	tvals, err = util.IntersectionTValsP(b1, b2, pt1, pt2)
	if err != nil {
		return nil
	}
	t = tvals[0]
	if t > 0 && t < 1 {
		ont := 1 - t
		tmp := []float64{ont*b1[0] + t*b2[0], ont*b1[1] + t*b2[1]}
		return Line(rp1, tmp)
	}

	return nil
}
