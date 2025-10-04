package graphics2d

// Generalized version of StrokeProc to replace explicit TraceProcs in StrokeProc
// RHS and LHS path processors
// Post side processing path processor
// if closed => two closed paths
// if open => single closed path with end caps
// Cap types - butt [default], round, square

// NStrokeProc defines the default width, side path processors and cap types of the stroke.
type NStrokeProc struct {
	RHSProc      PathProcessor
	LHSProc      PathProcessor
	PostSideProc PathProcessor
	Width        float64 // Default
	// point(pt, r) []part
	PointFunc func([]float64, float64) [][][]float64
	// cap(part1, pt, part2) []part
	CapStartFunc func([][]float64, []float64, [][]float64) [][][]float64
	CapEndFunc   func([][]float64, []float64, [][]float64) [][][]float64
}

// Process implements the PathProcessor interface and will return either one or two paths
// depending on whether the path is open or closed.
func (sp *NStrokeProc) Process(p *Path) []*Path {
	// Points are their own special case
	steps := p.Steps()
	if len(steps) == 1 {
		np := PartsToPath(sp.PointFunc(steps[0][0], sp.Width)...)
		np.Close()
		return []*Path{np}
	}

	// Calculate paths for each side - want the LHS one backwards
	rhs := sp.RHSProc.Process(p)[0]
	lhs := sp.LHSProc.Process(p)[0]

	if sp.PostSideProc != nil {
		// Run the post trace path processor on both side paths
		rhs = rhs.Process(sp.PostSideProc)[0]
		lhs = lhs.Process(sp.PostSideProc)[0]
	}

	// Need to switch the direction of the LHS path
	lhs = lhs.Reverse()

	if p.closed {
		// Process has already performed the last join
		rhs.Close()
		lhs.Close()
		return []*Path{rhs, lhs}
	}

	// Path is open, construct end caps and concatenate RHS with LHS, return it
	rhsp, lhsp := rhs.Parts(), lhs.Parts()

	both := make([][][]float64, 0, len(rhsp)+len(lhsp)+2)
	both = append(both, rhsp...)

	rhsl := rhsp[len(rhsp)-1]
	// cap pt is centroid of e1 and s2
	x := (rhsl[len(rhsl)-1][0] + lhsp[0][0][0]) / 2
	y := (rhsl[len(rhsl)-1][1] + lhsp[0][0][1]) / 2
	pt := []float64{x, y}
	both = append(both, sp.CapEndFunc(rhsl, pt, lhsp[0])...)

	both = append(both, lhsp...)

	lhsl := lhsp[len(lhsp)-1]
	// cap pt is centroid of e1 and s2
	x = (lhsl[len(lhsl)-1][0] + rhsp[0][0][0]) / 2
	y = (lhsl[len(lhsl)-1][1] + rhsp[0][0][1]) / 2
	pt = []float64{x, y}
	both = append(both, sp.CapStartFunc(lhsl, pt, rhsp[0])...)

	bp := PartsToPath(both...)
	bp.Close()
	return []*Path{bp}
}
