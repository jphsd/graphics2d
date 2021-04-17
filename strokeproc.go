package graphics2d

import (
	"math"

	. "github.com/jphsd/graphics2d/util"
)

// Constant width path stroker.
// if closed => two closed paths
// if open => single closed path with end caps
// Cap types - butt [default], round, square
// Join types - round, bevel [default], miter

// StrokeProc defines the width, join and cap types of the stroke.
type StrokeProc struct {
	Width        float64
	hw           float64 // hw - half width
	PointFunc    func([]float64, float64) [][][]float64
	JoinFunc     func([]float64, []float64, []float64) [][][]float64
	CapFunc      func([]float64, []float64, []float64) [][][]float64
	CapStartFunc func([]float64, []float64, []float64) [][][]float64
	CapEndFunc   func([]float64, []float64, []float64) [][][]float64
}

// NewStrokeProc creates a stroke path processor with width w, the bevel join and butt cap types.
func NewStrokeProc(w float64) *StrokeProc {
	if w < 0 {
		w = -w
	}
	return &StrokeProc{w, w / 2, PointCircle, JoinBevel, CapButt, nil, nil} // 10 degrees
}

// Process implements the PathProcessor interface and will return either one or two paths
// depending on whether the path is open or closed. Note that inside joins have a tie inside
// the stroke.
func (sp *StrokeProc) Process(p *Path) []*Path {
	steps := p.Steps()

	// Points are their own special case
	if len(steps) == 1 {
		np, _ := PartsToPath(sp.PointFunc(steps[0][0], sp.Width)...)
		np.Close()
		return []*Path{np}
	}

	// Preprocess curves into safe forms.
	p = p.Simplify()
	steps = p.Steps()

	// For each step, calculate the RHS start and end offsets [start/end][x/y/dx/dy]
	n := len(steps)
	var stepoffs, fpts [][][]float64
	if p.closed {
		stepoffs = make([][][]float64, n-1, n)
		fpts = make([][][]float64, n-1, n)
	} else {
		stepoffs = make([][][]float64, n-1)
		fpts = make([][][]float64, n-1)
	}
	cp := steps[0][0]
	for i := 1; i < n; i++ {
		pts := toPoints(cp, steps[i])
		tmp := make([][]float64, 2)
		tmp[0], tmp[1] = DeCasteljau(pts, 0), DeCasteljau(pts, 1)
		dx, dy := norm(tmp[0][3], -tmp[0][2])
		tmp[0][2], tmp[0][3] = dx*sp.hw, dy*sp.hw
		dx, dy = norm(tmp[1][3], -tmp[1][2])
		tmp[1][2], tmp[1][3] = dx*sp.hw, dy*sp.hw
		stepoffs[i-1] = tmp
		fpts[i-1] = pts
		cp = tmp[1]
	}
	if p.closed {
		pts := [][]float64{stepoffs[n-2][1][0:2], stepoffs[0][0][0:2]}
		if !EqualsP(pts[0], pts[1]) {
			tmp := make([][]float64, 2)
			tmp[0], tmp[1] = DeCasteljau(pts, 0), DeCasteljau(pts, 1)
			dx, dy := norm(tmp[0][3], -tmp[0][2])
			tmp[0][2], tmp[0][3] = dx*sp.hw, dy*sp.hw
			dx, dy = norm(tmp[1][3], -tmp[1][2])
			tmp[1][2], tmp[1][3] = dx*sp.hw, dy*sp.hw
			stepoffs = append(stepoffs, tmp)
			fpts = append(fpts, pts)
		}
	}

	// Calculate the RHS path by LineTransforming the steps and handling the joins
	n = len(stepoffs)
	rhs := make([][][]float64, n)
	for i := 0; i < n; i++ {
		tmp := stepoffs[i]
		xfm := LineTransform(tmp[0][0], tmp[0][1],
			tmp[1][0], tmp[1][1],
			tmp[0][0]+tmp[0][2], tmp[0][1]+tmp[0][3],
			tmp[1][0]+tmp[1][2], tmp[1][1]+tmp[1][3])
		rhs[i] = xfm.Apply(fpts[i]...)
	}

	// Compute the joins
	nrhs := make([][][]float64, 0, 2*n-1)
	nrhs = append(nrhs, rhs[0])
	for i := 1; i < n; i++ {
		last := rhs[i-1][len(rhs[i-1])-1]
		if !EqualsP(last, rhs[i][0]) {
			nrhs = append(nrhs, sp.JoinFunc(last, stepoffs[i][0], rhs[i][0])...)
		}
		nrhs = append(nrhs, rhs[i])
	}
	rhs = nrhs

	// Flip the steps and s/e offsets
	stepoffs = reverseOffs(stepoffs)
	fpts = ReverseParts(fpts)

	// Calculate the LHS path by LineTransforming the steps and handling the joins
	lhs := make([][][]float64, n)
	for i := 0; i < n; i++ {
		tmp := stepoffs[i]
		xfm := LineTransform(tmp[0][0], tmp[0][1],
			tmp[1][0], tmp[1][1],
			tmp[0][0]+tmp[0][2], tmp[0][1]+tmp[0][3],
			tmp[1][0]+tmp[1][2], tmp[1][1]+tmp[1][3])
		lhs[i] = xfm.Apply(fpts[i]...)
	}

	// Compute the joins
	nlhs := make([][][]float64, 0, 2*n-1)
	nlhs = append(nlhs, lhs[0])
	for i := 1; i < n; i++ {
		last := lhs[i-1][len(lhs[i-1])-1]
		if !EqualsP(last, lhs[i][0]) {
			nlhs = append(nlhs, sp.JoinFunc(last, stepoffs[i][0], lhs[i][0])...)
		}
		nlhs = append(nlhs, lhs[i])
	}
	lhs = nlhs

	var res []*Path
	if p.closed {
		// Close the RHS and LHS paths and return them
		rhsl := rhs[len(rhs)-1]
		if !EqualsP(rhsl[len(rhsl)-1], rhs[0][0]) {
			rhs = append(rhs, sp.JoinFunc(rhsl[len(rhsl)-1], stepoffs[0][0], rhs[0][0])...)
		}
		lhsl := lhs[len(lhs)-1]
		if !EqualsP(lhsl[len(lhsl)-1], lhs[0][0]) {
			lhs = append(lhs, sp.JoinFunc(lhsl[len(lhsl)-1], stepoffs[0][0], lhs[0][0])...)
		}
		rhsp, _ := PartsToPath(rhs...)
		rhsp.Close()
		lhsp, _ := PartsToPath(lhs...)
		lhsp.Close()
		res = []*Path{rhsp, lhsp}
	} else {
		if sp.CapEndFunc == nil {
			sp.CapEndFunc = sp.CapFunc
		}
		if sp.CapStartFunc == nil {
			sp.CapStartFunc = sp.CapFunc
		}
		// Path is open, construct end caps and concatenate RHS with LHS, return it
		both := make([][][]float64, 0, len(rhs)+len(lhs)+2)
		both = append(both, rhs...)
		rhsl := rhs[len(rhs)-1]
		both = append(both, sp.CapEndFunc(rhsl[len(rhsl)-1], stepoffs[0][0], lhs[0][0])...)
		both = append(both, lhs...)
		lhsl := lhs[len(lhs)-1]
		both = append(both, sp.CapStartFunc(lhsl[len(lhsl)-1], stepoffs[len(stepoffs)-1][1], rhs[0][0])...)
		bp, _ := PartsToPath(both...)
		bp.Close()
		res = []*Path{bp}
	}

	return res
}

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
	return MakeArcParts(pt[0], pt[1], w/2, 0, math.Pi*2)
}

// norm converts a normal to a unit normal
func norm(dx, dy float64) (float64, float64) {
	d := math.Sqrt(dx*dx + dy*dy)
	if Equals(0, d) {
		return 0, 0
	}
	return dx / d, dy / d
}

// [part][start/end][x/y/dx/dy]
func reverseOffs(parts [][][]float64) [][][]float64 {
	n := len(parts)
	res := make([][][]float64, n)
	for i, j := 0, n-1; i < n; i++ {
		res[i] = make([][]float64, 2)
		res[i][0], res[i][1] = parts[j][1], parts[j][0]
		// flip dx and dy
		res[i][0][2], res[i][0][3] = -res[i][0][2], -res[i][0][3]
		res[i][1][2], res[i][1][3] = -res[i][1][2], -res[i][1][3]
		j--
	}
	return res
}

// JoinBevel creates a bevel join from e1 to s2.
func JoinBevel(e1, p, s2 []float64) [][][]float64 {
	return [][][]float64{{e1, s2}}
}

// JoinRound creates a round join from e1 to s2, centered on p.
func JoinRound(e1, p, s2 []float64) [][][]float64 {
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	a1 := math.Atan2(dy, dx)
	a2 := LineAngle(p, s2)
	da := a2 - a1
	if da < -math.Pi {
		da += 2 * math.Pi
	} else if da > math.Pi {
		da -= 2 * math.Pi
	}
	if da < 0 {
		// inside angle
		return [][][]float64{{e1, s2}}
	}
	r := math.Sqrt(dx*dx + dy*dy)
	return MakeArcParts(p[0], p[1], r, a1, da)
}

// MiterJoin describes the limit and alternative function to use when the limit is exceeded
// for a miter join.
type MiterJoin struct {
	MiterLimit   float64
	MiterAltFunc func([]float64, []float64, []float64) [][][]float64
}

// NewMiterJoin creates a default MiterJoin with the limit set to 10 degrees and the alternative
// function to JoinBevel.
func NewMiterJoin() *MiterJoin {
	return &MiterJoin{10 * math.Pi / 180, JoinBevel}
}

// JoinMiter creates a miter join from e1 to s2 unless the moter limit is exceeded in which
// case the alternative function is used to perform the join.
func (mj *MiterJoin) JoinMiter(e1, p, s2 []float64) [][][]float64 {
	dx1, dy1 := e1[0]-p[0], e1[1]-p[1]
	a1 := math.Atan2(dy1, dx1)
	dx2, dy2 := s2[0]-p[0], s2[1]-p[1]
	a2 := math.Atan2(dy2, dx2)
	da := a2 - a1
	if da < -math.Pi {
		da += 2 * math.Pi
	} else if da > math.Pi {
		da -= 2 * math.Pi
	}
	if da < 0 {
		// inside angle
		return JoinBevel(e1, p, s2)
	}
	if da > math.Pi-mj.MiterLimit {
		// miter limit exceeded
		if mj.MiterAltFunc != nil {
			return mj.MiterAltFunc(e1, p, s2)
		}
		return JoinBevel(e1, p, s2)
	}
	// tangent -dy, dx
	ts, err := IntersectionTVals(e1[0], e1[1], e1[0]-dy1, e1[1]+dx1,
		s2[0], s2[1], s2[0]-dy2, s2[1]+dx2)
	if err != nil {
		return JoinBevel(e1, p, s2)
	}
	px := Lerp(ts[0], e1[0], e1[0]-dy1)
	py := Lerp(ts[0], e1[1], e1[1]+dx1)
	j := []float64{px, py}
	return [][][]float64{{e1, j}, {j, s2}}
}

// TODO - ArcJoin and MiterLimitJoin per
// https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-linejoin

// CapButt draws a line from e1 to s1.
func CapButt(e1, p, s1 []float64) [][][]float64 {
	return [][][]float64{{e1, s1}}
}

// CapRound draws a semicircle from e1 to s1 centered on p.
func CapRound(e1, p, s1 []float64) [][][]float64 {
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	offs := math.Atan2(dy, dx)
	r := math.Sqrt(dx*dx + dy*dy)
	return MakeArcParts(p[0], p[1], r, offs, math.Pi)
}

// CapSquare draws an extended square (stroke width/2) from e1 and s1.
func CapSquare(e1, p, s1 []float64) [][][]float64 {
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	return [][][]float64{{e1, e2}, {e2, s2}, {s2, s1}}
}

// CapHead draws an extended arrow head from e1 to extended p and then to s1.
func CapHead(e1, p, s1 []float64) [][][]float64 {
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	m := []float64{p[0] - dy, p[1] + dx}
	return [][][]float64{{e1, m}, {m, s1}}
}

// CapTail draws an extended arrow tail (stroke width/2) from extended e1 to p and then to extended  s1.
func CapTail(e1, p, s1 []float64) [][][]float64 {
	dx, dy := e1[0]-p[0], e1[1]-p[1]
	e2 := []float64{e1[0] - dy, e1[1] + dx}
	s2 := []float64{s1[0] - dy, s1[1] + dx}
	return [][][]float64{{e1, e2}, {e2, p}, {p, s2}, {s2, s1}}
}

// CapInvRound extends e1 and s1 and draws a semicircle that passes through p.
func CapInvRound(e1, p, s1 []float64) [][][]float64 {
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
