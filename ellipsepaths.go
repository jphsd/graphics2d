package graphics2d

import (
	"fmt"
	"math"

	"github.com/jphsd/graphics2d/util"
)

// A collection of elliptical path creation functions.

// Ellipse returns a closed path describing an ellipse with rx and ry rotated by xang from the x axis.
func Ellipse(c []float64, rx, ry, xang float64) *Path {
	ax, ay := c[0], c[1]
	np := PartsToPath(MakeArcParts(ax, ay, rx, 0, TwoPi)...)
	np.Close()
	// Reverse order
	xfm := Translate(ax, ay)
	xfm.Rotate(xang)
	xfm.Scale(1, ry/rx)
	xfm.Translate(-ax, -ay)
	return np.Process(xfm)[0]
}

// EllipseFromPoints returns a path describing the smallest ellipse containing points p1 and p2.
// If p1, p2 and c are colinear and not equidistant then nil is returned.
func EllipseFromPoints(p1, p2, c []float64) *Path {
	d1, d2 := util.DistanceESquared(p1, c), util.DistanceESquared(p2, c)
	if util.Collinear(p1, p2, c) {
		if !util.Equals(d1, d2) {
			return nil
		}
		// Points are on a circle
		return Circle(c, math.Sqrt(d1))
	}

	swap := d1 < d2
	if swap {
		p1, p2 = p2, p1
	}

	xang := util.LineAngle(c, p1)

	// Transform p1, p2 and c so that c is at the origin and p1 lies on the x axis.
	// Figure rx and ry from the transformed points.
	xfm := Rotate(-xang)
	xfm.Translate(-c[0], -c[1])
	pts := xfm.Apply(p1, p2, c)
	rx := pts[0][0]
	if rx < 0 {
		rx = -rx
	}
	rx2 := rx * rx
	// From x^2/rx^2 + y^2/ry^2 = 1
	ry2 := rx2 * pts[1][1] * pts[1][1] / (rx2 - pts[1][0]*pts[1][0])
	ry := math.Sqrt(ry2)

	return Ellipse(c, rx, ry, xang)
}

// EllipticalArc returns a path describing an arc starting at offs and ending at offs+ang on the ellipse
// defined by rx and ry rotated by xang from the x axis.
func EllipticalArc(c []float64, rx, ry, offs, ang, xang float64, s ArcStyle) *Path {
	// Angle conversion from elliptical space to circular: offs -> toffs, ang -> tang,
	// by generating points on unit circle for start and end, transforming them and finding
	// the new angles.
	xfm := Scale(1, rx/ry) // vs inverting the other

	offs -= xang
	sx, sy := math.Cos(offs), math.Sin(offs)
	end := offs + ang
	ex, ey := math.Cos(end), math.Sin(end)
	inv := xfm.Apply([]float64{sx, sy}, []float64{ex, ey})

	toffs, tend := math.Atan2(inv[0][1], inv[0][0]), math.Atan2(inv[1][1], inv[1][0])
	tang := tend - toffs
	if ang < 0 && tang > 0 {
		tang -= TwoPi
	} else if ang > 0 && tang < 0 {
		tang += TwoPi
	}

	np := PartsToPath(MakeArcParts(0, 0, rx, toffs, tang)...)

	switch s {
	case ArcChord:
		np.Close()
	case ArcPie:
		np.AddStep(c)
		np.Close()
	}

	// Turn circle into rotated ellipse
	xfm = Translate(c[0], c[1])
	xfm.Rotate(xang)
	xfm.Scale(1, ry/rx)

	return np.Process(xfm)[0]
}

// EllipticalArcFromPoint returns a path describing an ellipse arc from a point. The ratio of rx to ry
// is specified by rxy.
func EllipticalArcFromPoint(pt, c []float64, rxy, ang, xang float64, s ArcStyle) *Path {
	dx := pt[0] - c[0]
	dy := pt[1] - c[1]
	offs := math.Atan2(dy, dx)

	// calc rx and ry from pt and rxy
	xfm := Rotate(-xang)
	xfm.Translate(-c[0], -c[1])
	pts := xfm.Apply(pt)
	dx, dy = pts[0][0], pts[0][1]
	ry := math.Hypot(dx/rxy, dy)
	rx := rxy * ry

	return EllipticalArc(c, rx, ry, offs, ang, xang, s)
}

// EllipticalArcFromPoints returns a path describing the smallest ellipse arc from a point p1 to p2 (ccw).
// If p1, p2 and c are colinear and not equidistant then nil is returned.
func EllipticalArcFromPoints(p1, p2, c []float64, s ArcStyle) *Path {
	d1, d2 := util.DistanceESquared(p1, c), util.DistanceESquared(p2, c)
	if util.Collinear(p1, p2, c) {
		if !util.Equals(d1, d2) {
			// No solution
			return nil
		}
		// Infinite solutions - choose circular
		return ArcFromPoint(p1, c, Pi, s)
	}

	p1a := util.LineAngle(c, p1)
	p2a := util.LineAngle(c, p2)

	offs := p1a
	xang := p1a
	swap := d1 < d2
	if swap {
		xang = p2a
		p1, p2 = p2, p1
	}

	// Map [-pi,pi] to [0,2pi]
	if p1a < 0 {
		p1a = TwoPi + p1a
	}
	if p2a < 0 {
		p2a = TwoPi + p2a
	}
	var ang float64
	if p1a < p2a {
		ang = p2a - p1a
	} else {
		ang = TwoPi - p1a + p2a
	}

	// Transform p1, p2 and c so that c is at the origin and p1 lies on the x axis
	xfm := Rotate(-xang)
	xfm.Translate(-c[0], -c[1])
	pts := xfm.Apply(p1, p2, c)
	rx := pts[0][0]
	if rx < 0 {
		rx = -rx
	}
	rx2 := rx * rx
	// From x^2/rx^2 + y^2/ry^2 = 1
	ry2 := rx2 * pts[1][1] * pts[1][1] / (rx2 - pts[1][0]*pts[1][0])
	ry := math.Sqrt(ry2)

	return EllipticalArc(c, rx, ry, offs, ang, xang, s)
}

// EllipticalArcFromPoints2 provides a specification similar to that found in the SVG11 standard
// where the center is calculated from the given rx and ry values and flag specifications.
// See https://www.w3.org/TR/SVG11/implnote.html#ArcImplementationNotes
func EllipticalArcFromPoints2(p1, p2 []float64, rx, ry, xang float64, arc, swp bool, s ArcStyle) *Path {
	if util.Equals(rx, 0) || util.Equals(ry, 0) || util.EqualsP(p1, p2) {
		return Line(p1, p2)
	}

	// F6.5 Step1
	x1 := math.Cos(xang)*(p1[0]-p2[0])/2 + math.Sin(xang)*(p1[1]-p2[1])/2
	y1 := -math.Sin(xang)*(p1[0]-p2[0])/2 + math.Cos(xang)*(p1[1]-p2[1])/2

	// Fix rx and ry if necs (see F6.6)
	if rx < 0 {
		rx = -rx
	}
	if ry < 0 {
		ry = -ry
	}
	x12, y12, rx2, ry2 := x1*x1, y1*y1, rx*rx, ry*ry
	l := x12/rx2 + y12/ry2
	if l > 1 {
		// radii are too small
		sqrl := math.Sqrt(l)
		rx, ry = sqrl*rx, sqrl*ry
		rx2, ry2 = rx*rx, ry*ry
	}

	// F6.5 Step2
	val := (rx2*ry2 - rx2*y12 - ry2*x12) / (rx2*y12 + ry2*x12)
	if util.Equals(val, 0) {
		// Fix -0 before sqrt
		val = 0
	}
	tmp := math.Sqrt(val)
	if tmp != tmp {
		panic("constant is NaN")
	}
	if arc == swp {
		tmp = -tmp
	}
	cx, cy := tmp*rx*y1/ry, -tmp*ry*x1/rx

	// F6.5 Step3
	cx1 := math.Cos(xang)*cx - math.Sin(xang)*cy + (p1[0]+p2[0])/2
	cy1 := math.Sin(xang)*cx + math.Cos(xang)*cy + (p1[1]+p2[1])/2

	// F6.5 Step4 (differs from spec)
	c := []float64{cx1, cy1}
	offs := util.LineAngle(c, p1)
	ang, _, _ := util.AngleBetweenLines(c, p1, c, p2)

	if !swp && ang > 0 {
		ang -= TwoPi
	} else if swp && ang < 0 {
		ang += TwoPi
	}

	return EllipticalArc([]float64{cx1, cy1}, rx, ry, offs, ang, xang, s)
}

// IrregularEllipse uses different rx and ry values for each quadrant of an ellipse. disp (-1,1)
// determines how far along either rx1 (+ve) or rx2 (-ve), ry2 extends from (ry1 extends from c).
func IrregularEllipse(c []float64, rx1, rx2, ry1, ry2, disp, xang float64) *Path {
	// Limit disp
	if disp > 1 {
		disp = 1
	} else if disp < -1 {
		disp = -1
	}

	var dx float64
	if disp < 0 {
		dx = rx2 * disp
	} else {
		dx = rx1 * disp
	}
	rxx1 := rx1 - dx
	rxx2 := rx2 + dx

	ang := HalfPi
	var offs float64

	parts := []Part{}
	parts = append(parts, EllipticalArc([]float64{0, 0}, rx1, ry1, offs, ang, 0, ArcOpen).Parts()...)
	offs += ang
	parts = append(parts, EllipticalArc([]float64{0, 0}, rx2, ry1, offs, ang, 0, ArcOpen).Parts()...)
	offs += ang
	parts = append(parts, EllipticalArc([]float64{dx, 0}, rxx2, ry2, offs, ang, 0, ArcOpen).Parts()...)
	offs += ang
	parts = append(parts, EllipticalArc([]float64{dx, 0}, rxx1, ry2, offs, ang, 0, ArcOpen).Parts()...)

	// Move to c and rotate by xang
	xfm := Translate(c[0], c[1])
	xfm.Rotate(xang)

	path := PartsToPath(parts...).Process(xfm)[0]
	path.Close()
	return path
}

// Egg uses IrregularEllipse to generate an egg shape with the specified width and height. The waist is
// specified as a percentage distance along the height axis (from the base). The egg is rotated by xang.
func Egg(c []float64, w, h, d, xang float64) *Path {
	rx := w / 2
	ryb := h * d
	ryt := h - ryb
	return IrregularEllipse(c, rx, rx, ryt, ryb, 0, xang)
}

// RightEgg uses IrregularEllipse to generate an egg shape with the specified height and a semicircular base.
// The waist is specified as a percentage distance along the height axis (from the base). The egg is rotated by xang.
func RightEgg(c []float64, h, d, xang float64) *Path {
	ryb := h * d
	ryt := h - ryb
	return IrregularEllipse(c, ryb, ryb, ryt, ryb, 0, xang)
}

// From Robert Dixon, Mathographics, 1987
//   Moss: d = 1 / (4 - sqr(2))
//   Golden: d = 1 / (1 + phi) [phi = (1 + sqr(5)) / 2]
//   Cundy Rollett: d = 2 / (3 + sqr(3))
//   Thom1: d = 7 / 16
//   Thom2: d = 3 / (2 + sqr(10))
// Stephanie Moss ???
// Martyn Cundy and AP Rollett, Mathematical Models, 1951
// Alexander Thom, Megalithic Sites in Britain, 1967
// Note - above are for polycentric curves, not true ellipses.

// Nellipse takes a slice of ordered foci, assumed to be on the hull of a convex polygon, and a length, and uses them
// to construct a Gardener's ellipse (an approximation that ignores foci within the hull). The closed path will be
// made up of twice as many arcs as there are foci. If the length isn't sufficient to wrap the foci, then nil is
// returned.
func Nellipse(l float64, foci ...[]float64) *Path {
	np := len(foci)

	// Trivial cases
	switch np {
	case 0:
		return nil
	case 1:
		// Circle
		return Circle(foci[0], l/2)
	case 2:
		// Single ellipse
		ml := util.DistanceE(foci[0], foci[1]) * 2
		if l < ml {
			return nil
		}
		fe := &fflEllipse{foci[0], foci[1], 0, 1, l, nil, nil, -1, -1, nil}
		c, rx, ry, th := fe.toEllipseArgs()
		return Ellipse(c, rx, ry, th)
	}

	// Check l wraps the foci
	ml := 0.0
	for i := 0; i < np-1; i++ {
		ml += util.DistanceE(foci[i], foci[i+1])
	}
	ml += util.DistanceE(foci[np-1], foci[0])
	if l < ml {
		return nil
	}

	// First set of ellipses - two (R & L) for each face tangent
	ellipsesR, ellipsesL := calcEllipses(foci, l, false)

	// Create a map of them
	em := make(map[string]*fflEllipse)
	for i, ellip := range ellipsesR {
		addEllipse(em, ellip)
		ellipsesR[i] = addEllipse(em, ellip)
	}
	for i, ellip := range ellipsesL {
		ellipsesL[i] = addEllipse(em, ellip)
	}

	// Reverse points and calc the second set (ie reversed face tangents)
	rellipsesR, rellipsesL := calcEllipses(ReversePoints(foci), l, true)

	// Fix point reversal in ellipses and swap intercepts
	for i, ellip := range rellipsesR {
		ellip.f1i = np - 1 - ellip.f1i
		ellip.f2i = np - 1 - ellip.f2i
		ellip.i1, ellip.i2 = ellip.i2, ellip.i1
		ellip.i1i, ellip.i2i = ellip.i2i, ellip.i1i
		ellip.i2i += np
		rellipsesR[i] = addEllipse(em, ellip)
	}
	for i, ellip := range rellipsesL {
		ellip.f1i = np - 1 - ellip.f1i
		ellip.f2i = np - 1 - ellip.f2i
		ellip.i1, ellip.i2 = ellip.i2, ellip.i1
		ellip.i1i, ellip.i2i = ellip.i2i, ellip.i1i
		ellip.i1i += np
		rellipsesL[i] = addEllipse(em, ellip)
	}

	// Ellipse intercepts are complete, construct arcs
	i2e := make([]*fflEllipse, len(em))
	for _, v := range em {
		i2e[v.i1i] = v
		c, rx, ry, offs, ang, th := v.toEllipseArcArgs()
		v.path = EllipticalArc(c, rx, ry, offs, ang, th, ArcOpen)
	}

	// Construct final path from individual arcs
	ce := ellipsesR[0]
	ni := ce.i2i
	path := ce.path
	for i := 1; i < 2*np; i++ {
		ce = i2e[ni]
		path.Concatenate(ce.path)
		ni = ce.i2i
	}
	path.Close()

	return path
}

func calcEllipses(points [][]float64, l float64, flip bool) ([]*fflEllipse, []*fflEllipse) {
	np := len(points)
	resR := make([]*fflEllipse, np)
	resL := make([]*fflEllipse, np)
	q := points[:]
	for i := range np {
		ellipseR, ellipseL := ellipseCalc(q, l, i, flip)
		// Adjust indices to handle rotation
		ellipseR.f1i += i
		if ellipseR.f1i >= np {
			ellipseR.f1i -= np
		}
		ellipseR.f2i += i
		if ellipseR.f2i >= np {
			ellipseR.f2i -= np
		}
		resR[i] = ellipseR
		ellipseL.f1i += i
		if ellipseL.f1i >= np {
			ellipseL.f1i -= np
		}
		ellipseL.f2i += i
		if ellipseL.f2i >= np {
			ellipseL.f2i -= np
		}
		resL[i] = ellipseL
		// Rotate points
		p := q[0]
		q = q[1:]
		q = append(q, p)
	}
	return resR, resL
}

// Calculate the ellipse for p0 and l. Returns f1 (= foci[0]), f2 and lr for the R ellipse
// and f0 (= foci[1]), f2 and ll for the L ellipse. Flip controls the side of line calculation
func ellipseCalc(foci [][]float64, l float64, id int, flip bool) (*fflEllipse, *fflEllipse) {
	f0 := foci[1] // f0 and f1 defines the face tangent
	f1 := foci[0]
	fi := findFurthest(foci)

	// Skip foci < fi
	for i := 0; i < fi-1; i++ {
		l -= util.DistanceE(foci[i], foci[i+1])
	}

	// Test points fi thru n until we find one that has no poly points on the RHS of it
	np := len(foci)
	for i := fi; i < np; i++ {
		cp := foci[i]
		l -= util.DistanceE(foci[i-1], cp)

		pt, d := findEllipseIntersection(f0, f1, cp, l)
		if i == np-1 {
			// No need to test
			return &fflEllipse{f1, cp, 0, i, l + d, pt, nil, id, -1, nil},
				&fflEllipse{f0, cp, 1, i, l + util.DistanceE(f0, f1) + util.DistanceE(f0, cp), nil, pt, -1, id, nil}
		} else {
			sol := util.SideOfLine(cp, pt, foci[i+1])
			if (flip && sol > 0) || (!flip && sol < 0) {
				return &fflEllipse{f1, cp, 0, i, l + d, pt, nil, id, -1, nil},
					&fflEllipse{f0, cp, 1, i, l + util.DistanceE(f0, f1) + util.DistanceE(f0, cp), nil, pt, -1, id, nil}
			}
		}
	}

	// Should never get here!
	panic("ellipseCalc failed")
}

// Assumes points are in adjacency order and for the first point returns the index of the point furthest away from it
// Assumes #points > 2
func findFurthest(points [][]float64) int {
	np := len(points)
	max := -1.0
	cur := -1
	for i := 2; i < np; i++ {
		d := util.DistanceE(points[0], points[i])
		if d > max {
			max = d
			cur = i
		}
	}
	return cur
}

// Find p on line f0->f1, such that d(p,fi)+d(p,f1) = l
func findEllipseIntersection(f0, f1, fi []float64, l float64) ([]float64, float64) {
	a := util.DistanceE(f1, fi)
	th, _, _ := util.AngleBetweenLines(f1, f0, f1, fi)
	if th < 0 {
		th = -th
	}
	th = Pi - th

	// Use the Cosine rule to find b given a, th and l=b+c
	b := (l*l - a*a) / (2 * (l - a*math.Cos(th)))

	// Get unit vector for f0->f1 and scale it by b
	vec := util.VecNormalize(util.Vec(f0, f1))
	vec[0] *= b
	vec[1] *= b

	pt := []float64{f1[0] + vec[0], f1[1] + vec[1]}
	return pt, a
}

func addEllipse(em map[string]*fflEllipse, ellip *fflEllipse) *fflEllipse {
	h := ellip.hash()
	me := em[h]
	if me == nil {
		em[h] = ellip
	} else {
		// Merge the intercepts
		if me.i1i == -1 {
			me.i1, me.i1i = ellip.i1, ellip.i1i
		} else {
			me.i2, me.i2i = ellip.i2, ellip.i2i
		}
	}
	return me
}

// fflEllipse is an ellipse defined by it's two foci and the length of sides of the triangle formed
// by f1, f2 and a point on the circumference.
type fflEllipse struct {
	f1, f2   []float64 // foci
	f1i, f2i int       // index of foci
	l        float64   // length
	i1, i2   []float64 // face tangent intercepts
	i1i, i2i int       // index of face tangent intercepts
	path     *Path     // the arc path
}

// f1, f2, l => c, rx, ry, th
func (fe *fflEllipse) toEllipseArgs() ([]float64, float64, float64, float64) {
	c := []float64{(fe.f1[0] + fe.f2[0]) / 2, (fe.f1[1] + fe.f2[1]) / 2}
	dx, dy := fe.f2[0]-fe.f1[0], fe.f2[1]-fe.f1[1]
	df := math.Hypot(dx, dy)
	if fe.l < 2*df {
		return nil, 0, 0, 0
	}
	rx := (fe.l - df) / 2
	ry := math.Sqrt(rx*rx - df*df/4)
	theta := math.Atan2(dy, dx)
	return c, rx, ry, theta
}

// f1, f2, l, i1, i2 => c, rx, ry, offs, ang, th
func (fe *fflEllipse) toEllipseArcArgs() ([]float64, float64, float64, float64, float64, float64) {
	c := []float64{(fe.f1[0] + fe.f2[0]) / 2, (fe.f1[1] + fe.f2[1]) / 2}
	dx, dy := fe.f2[0]-fe.f1[0], fe.f2[1]-fe.f1[1]
	df := math.Hypot(dx, dy)
	if fe.l < 2*df {
		return nil, 0, 0, 0, 0, 0
	}
	rx := (fe.l - df) / 2
	ry := math.Sqrt(rx*rx - df*df/4)
	theta := math.Atan2(dy, dx)
	offs := util.LineAngle(c, fe.i1)
	ang, _, _ := util.AngleBetweenLines(c, fe.i1, c, fe.i2)
	return c, rx, ry, offs, ang, theta
}

func (fe *fflEllipse) hash() string {
	a, b := fe.f1i, fe.f2i
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("%d %d %f", a, b, fe.l)
}
