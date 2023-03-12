package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// A collection of elliptical path creation functions.

// Ellipse returns a closed path describing an ellipse with rx and ry rotated by xang from the x axis.
func Ellipse(c []float64, rx, ry, xang float64) *Path {
	ax, ay := c[0], c[1]
	np := PartsToPath(MakeArcParts(ax, ay, rx, 0, TwoPi)...)
	np.Close()
	xfm := NewAff3()
	// Reverse order
	xfm.Translate(ax, ay)
	xfm.Rotate(xang)
	xfm.Scale(1, ry/rx)
	xfm.Translate(-ax, -ay)
	return np.Transform(xfm)
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
	xfm := NewAff3()
	xfm.Rotate(-xang)
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
	xfm := NewAff3()
	xfm.Scale(1, rx/ry) // vs inverting the other

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
	xfm = NewAff3()
	xfm.Translate(c[0], c[1])
	xfm.Rotate(xang)
	xfm.Scale(1, ry/rx)

	return np.Transform(xfm)
}

// EllipticalArcFromPoint returns a path describing an ellipse arc from a point. The ratio of rx to ry
// is specified by rxy.
func EllipticalArcFromPoint(pt, c []float64, rxy, ang, xang float64, s ArcStyle) *Path {
	dx := pt[0] - c[0]
	dy := pt[1] - c[1]
	offs := math.Atan2(dy, dx)

	// calc rx and ry from pt and rxy
	xfm := NewAff3()
	xfm.Rotate(-xang)
	xfm.Translate(-c[0], -c[1])
	pts := xfm.Apply(pt)
	dx, dy = pts[0][0], pts[0][1]
	ry := math.Sqrt(dx*dx/(rxy*rxy) + dy*dy)
	rx := rxy * ry

	return EllipticalArc(c, rx, ry, offs, ang, xang, s)
}

// EllipticalArcFromPoints returns a path describing the smallest ellipse arc from a point p1 to p2 (ccw).
// If p1, p2 and c are colinear and not equidistant then nil is returned.
func EllipticalArcFromPoints(p1, p2, c []float64, s ArcStyle) *Path {
	d1, d2 := util.DistanceESquared(p1, c), util.DistanceESquared(p2, c)
	if util.Collinear(p1, p2, c) && !util.Equals(d1, d2) {
		// No solution
		return nil
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
	xfm := NewAff3()
	xfm.Rotate(-xang)
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

	ang := math.Pi / 2
	var offs float64

	parts := [][][]float64{}
	parts = append(parts, EllipticalArc([]float64{0, 0}, rx1, ry1, offs, ang, 0, ArcOpen).Parts()...)
	offs += ang
	parts = append(parts, EllipticalArc([]float64{0, 0}, rx2, ry1, offs, ang, 0, ArcOpen).Parts()...)
	offs += ang
	parts = append(parts, EllipticalArc([]float64{dx, 0}, rxx2, ry2, offs, ang, 0, ArcOpen).Parts()...)
	offs += ang
	parts = append(parts, EllipticalArc([]float64{dx, 0}, rxx1, ry2, offs, ang, 0, ArcOpen).Parts()...)

	// Move to c and rotate by xang
	xfm := NewAff3()
	xfm.Translate(c[0], c[1])
	xfm.Rotate(xang)

	path := PartsToPath(parts...).Transform(xfm)
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
