package graphics2d

import (
	"math"

	. "github.com/jphsd/graphics2d/util"
)

// A collection of part and path creation functions.

// MakeArcParts creates at least one cubic bezier that describes a curve from offs to
// offs+ang centered on {cx, cy} with radius r.
func MakeArcParts(cx, cy, r, offs, ang float64) [][][]float64 {
	n := 1
	a := ang
	rev := ang < 0
	if rev {
		a = -a
	}

	// Calculate number of curves to create - necessary since curve errors
	// are apparent for angles > Pi/2
	for true {
		if a < math.Pi/2 {
			break
		}
		a /= 2
		n *= 2
	}
	cp := CalcPointsForArc(a)
	if rev {
		cp = [][]float64{cp[3], cp[2], cp[1], cp[0]}
	}
	if rev {
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

const (
	ArcOpen ArcStyle = iota
	ArcChord
	ArcPie
)

// Arc returns a path with an arc centered on c with radius r from offs in the direction and length of ang.
func Arc(c []float64, r, offs, ang float64, s ArcStyle) *Path {
	parts := MakeArcParts(c[0], c[1], r, offs, ang)
	np, _ := PartsToPath(parts...)
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
	r := math.Sqrt(dx*dx + dy*dy)
	offs := math.Atan2(dy, dx)
	return Arc(c, r, offs, ang, s)
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
		r := math.Sqrt(dx*dx + dy*dy)
		offs := math.Atan2(dy, dx)
		tmp := MakeArcParts(cs[i][0], cs[i][1], r, offs, angs[i])
		last := tmp[len(tmp)-1]
		cp = last[len(last)-1]
		parts = append(parts, tmp...)
	}

	res, _ := PartsToPath(parts...)
	return res
}

// Ellipse returns a closed path describing an ellipse with rx and ry rotated by xang from the x axis.
func Ellipse(c []float64, rx, ry, xang float64) *Path {
	ax, ay := c[0], c[1]
	np, _ := PartsToPath(MakeArcParts(ax, ay, rx, 0, 2*math.Pi)...)
	np.Close()
	xfm := NewAff3()
	// Reverse order
	xfm.Translate(ax, ay)
	xfm.Rotate(xang)
	xfm.Scale(1, ry/rx)
	xfm.Translate(-ax, -ay)
	return np.Transform(xfm)
}

// EllipticalArc returns a path describing an ellipse arc with rx and ry rotated by xang from the x axis.
func EllipticalArc(c []float64, rx, ry, offs, ang, xang float64, s ArcStyle) *Path {
	ax, ay := c[0], c[1]
	np, _ := PartsToPath(MakeArcParts(ax, ay, rx, offs-xang, ang)...)

	switch s {
	case ArcChord:
		np.Close()
	case ArcPie:
		np.AddStep(c)
		np.Close()
	}

	xfm := NewAff3()
	// Reverse order
	xfm.Translate(ax, ay)
	xfm.Rotate(xang)
	xfm.Scale(1, ry/rx)
	xfm.Translate(-ax, -ay)
	return np.Transform(xfm)
}

// EllipticalArcFromPoint returns a path describing an ellipse arc from a point. The ratio of rx to ry
// is specified by rxy.
func EllipticalArcFromPoint(pt, c []float64, rxy, ang, xang float64, s ArcStyle) *Path {
	dx := pt[0] - c[0]
	dy := pt[1] - c[1]
	offs := math.Atan2(dy, dx)

	// calc rx and ry from dx, dy and rxy
	ry := math.Sqrt(dx*dx/rxy*rxy + dy*dy)
	rx := rxy * ry

	return EllipticalArc(c, rx, ry, offs, ang, xang, s)
}

// RegularPolygon returns a closed path describing an n-sided polygon given the initial edge.
func RegularPolygon(pt1, pt2 []float64, n int) *Path {
	da := 2 * math.Pi / float64(n)
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
