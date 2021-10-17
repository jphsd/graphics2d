package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// Constant width path stroker which uses TraceProc to calculate the two sides for it.
// if closed => two closed paths
// if open => single closed path with end caps
// Cap types - butt [default], round, square
// Join types - round, bevel [default], miter

// StrokeProc defines the width, join and cap types of the stroke.
type StrokeProc struct {
	RTraceProc *TraceProc
	LTraceProc *TraceProc
	// (pt, r) []part
	PointFunc func([]float64, float64) [][][]float64
	// (part1, pt, part2) []part
	CapFunc      func([][]float64, []float64, [][]float64) [][][]float64
	CapStartFunc func([][]float64, []float64, [][]float64) [][][]float64
	CapEndFunc   func([][]float64, []float64, [][]float64) [][][]float64
}

// NewStrokeProc creates a stroke path processor with width w, the bevel join and butt cap types.
func NewStrokeProc(w float64) *StrokeProc {
	if w < 0 {
		w = -w
	}
	return NewStrokeProcExt(w/2, -w/2, JoinBevel, JoinBevel, 0.5) // 10 degrees
}

// NewStrokeProcExt creates a stroke path processor where the trace join types and widths are specified
// separately for each side of the stroke. This allows the stroke to be offset to the left or right
// of the path being processed.
func NewStrokeProcExt(rw, lw float64,
	rjf, ljf func([][]float64, []float64, [][]float64) [][][]float64, d float64) *StrokeProc {
	if rw < 0 {
		rw = -rw
	}
	if lw > 0 {
		lw = -lw
	}
	return &StrokeProc{&TraceProc{rw, d, rjf}, &TraceProc{lw, d, ljf}, PointCircle, CapButt, nil, nil}
}

// Process implements the PathProcessor interface and will return either one or two paths
// depending on whether the path is open or closed.
func (sp *StrokeProc) Process(p *Path) []*Path {
	// Points are their own special case
	steps := p.Steps()
	if len(steps) == 1 {
		w := sp.RTraceProc.Width - sp.LTraceProc.Width
		np, _ := PartsToPath(sp.PointFunc(steps[0][0], w)...)
		np.Close()
		return []*Path{np}
	}

	// Calculate traces for each side
	rhs := sp.RTraceProc.ProcessParts(p)
	lhs := sp.LTraceProc.ProcessParts(p)

	if p.closed {
		// ProcessParts has already performed the last join
		rhsp, _ := PartsToPath(rhs...)
		rhsp.Close()
		lhsp, _ := PartsToPath(lhs...)
		lhsp.Close()
		return []*Path{rhsp, lhsp}
	}

	// Path is open, construct end caps and concatenate RHS with LHS, return it
	if sp.CapEndFunc == nil {
		sp.CapEndFunc = sp.CapFunc
	}
	if sp.CapStartFunc == nil {
		sp.CapStartFunc = sp.CapFunc
	}

	both := make([][][]float64, 0, len(rhs)+len(lhs)+2)
	both = append(both, rhs...)

	rhsl := rhs[len(rhs)-1]
	// cap pt is centroid of e1 and s2
	x := (rhsl[len(rhsl)-1][0] + lhs[0][0][0]) / 2
	y := (rhsl[len(rhsl)-1][1] + lhs[0][0][1]) / 2
	pt := []float64{x, y}
	both = append(both, sp.CapEndFunc(rhsl, pt, lhs[0])...)

	both = append(both, lhs...)

	lhsl := lhs[len(lhs)-1]
	// cap pt is centroid of e1 and s2
	x = (lhsl[len(lhsl)-1][0] + rhs[0][0][0]) / 2
	y = (lhsl[len(lhsl)-1][1] + rhs[0][0][1]) / 2
	pt = []float64{x, y}
	both = append(both, sp.CapStartFunc(lhsl, pt, rhs[0])...)
	bp, _ := PartsToPath(both...)
	bp.Close()
	return []*Path{bp}
}

// unit converts a normal to a unit normal
func unit(dx, dy float64) (float64, float64) {
	d := math.Sqrt(dx*dx + dy*dy)
	if util.Equals(0, d) {
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
