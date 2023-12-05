package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// PointSquare renders points as squares aligned in x/y.
func PointSquare(pt []float64, w float64) [][][]float64 {
	hw := w / 2
	sx, sy := pt[0]-hw, pt[1]-hw
	res := make([][][]float64, 4)
	res[0] = [][]float64{{sx, sy}, {sx + w, sy}}
	res[1] = [][]float64{{sx + w, sy}, {sx + w, sy + w}}
	res[2] = [][]float64{{sx + w, sy + w}, {sx, sy + w}}
	res[3] = [][]float64{{sx, sy + w}, {sx, sy}}
	return res
}

// PointDiamond renders points as diamonds aligned in x/y.
func PointDiamond(pt []float64, w float64) [][][]float64 {
	hw := w / 2
	sx, sy := pt[0]-hw, pt[1]-hw
	res := make([][][]float64, 4)
	res[0] = [][]float64{{pt[0], sy}, {sx + w, pt[1]}}
	res[1] = [][]float64{{sx + w, pt[1]}, {pt[0], sy + w}}
	res[2] = [][]float64{{pt[0], sy + w}, {sx, pt[1]}}
	res[3] = [][]float64{{sx, pt[1]}, {pt[0], sy}}
	return res
}

// PointCircle renders points as circles/
func PointCircle(pt []float64, w float64) [][][]float64 {
	return MakeArcParts(pt[0], pt[1], w/2, 0, TwoPi)
}

// Joins take the two parts to be joined, p1 and p2, and some center point, p.

// JoinBevel creates a bevel join from e1 to s2.
func JoinBevel(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s2 := p1[len(p1)-1], p2[0]
	return [][][]float64{{e1, s2}}
}

// JoinRound creates a round join from e1 to s2, centered on p.
func JoinRound(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s2 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	a1 := math.Atan2(dy, dx)
	a2 := util.LineAngle(p, s2)
	da := a2 - a1
	if da < -math.Pi {
		da += TwoPi
	} else if da > math.Pi {
		da -= TwoPi
	}
	if da < 0 {
		// inside angle
		return [][][]float64{{e1, s2}}
	}
	r := math.Hypot(dx, dy)
	return MakeArcParts(p[0], p[1], r, a1, da)
}

// MiterJoin describes the limit and alternative function to use when the limit is exceeded
// for a miter join.
type MiterJoin struct {
	MiterLimit   float64
	MiterAltFunc func([][]float64, []float64, [][]float64) [][][]float64
}

// NewMiterJoin creates a default MiterJoin with the limit set to 10 degrees and the alternative
// function to JoinBevel.
func NewMiterJoin() *MiterJoin {
	return &MiterJoin{10 * math.Pi / 180, JoinBevel}
}

// JoinMiter creates a miter join from e1 to s2 unless the moter limit is exceeded in which
// case the alternative function is used to perform the join.
func (mj *MiterJoin) JoinMiter(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s2 := p1[len(p1)-1], p2[0]
	dx1, dy1 := e1[0]-p[0], e1[1]-p[1]
	a1 := math.Atan2(dy1, dx1)
	dx2, dy2 := s2[0]-p[0], s2[1]-p[1]
	a2 := math.Atan2(dy2, dx2)
	da := a2 - a1
	if da < -math.Pi {
		da += TwoPi
	} else if da > math.Pi {
		da -= TwoPi
	}
	if da < 0 {
		// inside angle
		return JoinBevel(p1, p, p2)
	}
	if da > math.Pi-mj.MiterLimit {
		// miter limit exceeded
		if mj.MiterAltFunc != nil {
			return mj.MiterAltFunc(p1, p, p2)
		}
		return JoinBevel(p1, p, p2)
	}
	// tangent -dy, dx
	ts, err := util.IntersectionTVals(e1[0], e1[1], e1[0]-dy1, e1[1]+dx1,
		s2[0], s2[1], s2[0]-dy2, s2[1]+dx2)
	if err != nil {
		return JoinBevel(p1, p, p2)
	}
	px := util.Lerp(ts[0], e1[0], e1[0]-dy1)
	py := util.Lerp(ts[0], e1[1], e1[1]+dx1)
	j := []float64{px, py}
	return [][][]float64{{e1, j}, {j, s2}}
}

// TODO - ArcJoin and MiterLimitJoin per
// https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-linejoin

// Caps take the two 'parallel' parts to be joined, p1 and p2, and some center point, p.

// CapButt draws a line from e1 to s1.
func CapButt(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	return [][][]float64{{e1, s1}}
}

// CapRound draws a semicircle from e1 to s1 centered on p.
func CapRound(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1 := p1[len(p1)-1]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	offs := math.Atan2(dy, dx)
	r := math.Sqrt(dx*dx + dy*dy)
	return MakeArcParts(p[0], p[1], r, offs, math.Pi)
}

// CapInvRound extends e1 and s1 and draws a semicircle that passes through p.
func CapInvRound(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	offs := math.Atan2(dy, dx)
	r := math.Sqrt(dx*dx + dy*dy)
	tp := MakeArcParts(p[0]-dy, p[1]+dx, r, offs, -math.Pi)
	res := make([][][]float64, 1, len(tp)+2)
	res[0] = [][]float64{e1, e2}
	res = append(res, tp...)
	res = append(res, [][]float64{s2, s1})
	return res
}

// CapSquare draws an extended square (stroke width/2) from e1 and s1.
func CapSquare(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	return [][][]float64{{e1, e2}, {e2, s2}, {s2, s1}}
}

// CapHead draws an extended arrow head from e1 to extended p and then to s1.
func CapHead(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	m := []float64{p[0] - dy, p[1] + dx}
	return [][][]float64{{e1, m}, {m, s1}}
}

// CapTail draws an extended arrow tail (stroke width/2) from extended e1 to p and then to extended  s1.
func CapTail(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	return [][][]float64{{e1, e2}, {e2, p}, {p, s2}, {s2, s1}}
}

// OvalCap contains the ratio of rx to ry for the oval and a center line offset
type OvalCap struct {
	Rxy  float64 // Ratio of Rx to Ry
	Offs float64 // Offset from center line [-1,1] -1 = LHS, 0 = centerline, 1 = RHS
}

// CapOval creates a half oval with ry = w/2 and rx = Rxy * ry
func (oc *OvalCap) CapOval(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1 := p1[len(p1)-1]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	offs := math.Atan2(dy, dx)
	ry := math.Sqrt(dx*dx + dy*dy)
	rx := ry * oc.Rxy

	if util.Equals(oc.Offs, 0) {
		return EllipticalArc(p, rx, ry, offs, math.Pi, offs-HalfPi, ArcOpen).Parts()
	}

	// Construct two quarter arcs with new rys and cp
	s1 := p2[0]
	t := (oc.Offs + 1) / 2
	omt := 1 - t
	cp := []float64{e1[0]*omt + s1[0]*t, e1[1]*omt + s1[1]*t}
	d := 2 * ry
	ry1, ry2 := d*t, d*omt
	res := EllipticalArc(cp, rx, ry1, offs, HalfPi, offs-HalfPi, ArcOpen).Parts()
	return append(res, EllipticalArc(cp, rx, ry2, offs+HalfPi, HalfPi, offs-HalfPi, ArcOpen).Parts()...)
}

// CapInvOval creates an inverted half oval with rx = w/2 and ry = rxy * rx
// Offs is ignored
func (oc *OvalCap) CapInvOval(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	offs := math.Atan2(dy, dx)
	xoffs := offs - HalfPi
	ry := math.Sqrt(dx*dx + dy*dy)
	rx := ry * oc.Rxy
	tp := EllipticalArc([]float64{p[0] - dy*oc.Rxy, p[1] + dx*oc.Rxy}, rx, ry, offs, -math.Pi, xoffs, ArcOpen).Parts()
	res := make([][][]float64, 1, len(tp)+2)
	res[0] = [][]float64{e1, e2}
	res = append(res, tp...)
	res = append(res, [][]float64{s2, s1})
	return res
}

// RSCap contains the percentage [0,1] of the corner taken up by an arc. Perc = 1 is equivalent to CapRound
// Perc = 0, to CapSquare.
type RSCap struct {
	Perc float64
}

// CapRoundedSquare creates a square cap with rounded corners
func (rc *RSCap) CapRoundedSquare(p1 [][]float64, p []float64, p2 [][]float64) [][][]float64 {
	e1, s1 := p1[len(p1)-1], p2[0]
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	r := math.Sqrt(dx*dx+dy*dy) * rc.Perc
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	parts := [][][]float64{{e1, e2}, {e2, s2}, {s2, s1}}
	path := PartsToPath(parts...)
	rp := &RoundedProc{r}
	return path.Process(rp)[0].Parts()
}
