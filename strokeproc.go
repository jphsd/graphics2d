package graphics2d

// Generalized version of StrokeProc to replace explicit TraceProcs
// RHS and LHS path processors
// Post side processing path processor
// if closed => two closed paths
// if open => single closed path with end caps
// Cap types - butt [default], round, square

// StrokeProc defines the default width, side path processors and cap types of the stroke.
type StrokeProc struct {
	RHSProc      PathProcessor
	LHSProc      PathProcessor
	PostSideProc PathProcessor
	Width        float64 // Default
	// point(pt, r) []part
	PointFunc func([]float64, float64) []Part
	// cap(part1, pt, part2) []part
	CapStartFunc func(Part, []float64, Part) []Part
	CapEndFunc   func(Part, []float64, Part) []Part
}

// Process implements the PathProcessor interface and will return either one or two paths
// depending on whether the path is open or closed.
func (sp *StrokeProc) Process(p *Path) []*Path {
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

	both := make([]Part, 0, len(rhsp)+len(lhsp)+2)
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

// NewStrokeProc creates a trace stroke path processor with width w, the bevel join and butt cap types.
func NewStrokeProc(w float64) *StrokeProc {
	if w < 0 {
		w = -w
	}
	return NewStrokeProcExt(w/2, -w/2, JoinBevel, 0.5, CapButt) // 10 degrees
}

// NewStrokeProcExt creates a trace stroke path processor where the widths are specified
// separately for each side of the stroke. This allows the stroke to be offset to the left or right
// of the path being processed.
func NewStrokeProcExt(rw, lw float64,
	jf func(Part, []float64, Part) []Part,
	d float64,
	cf func(Part, []float64, Part) []Part) *StrokeProc {
	if rw < 0 {
		rw = -rw
	}
	if lw > 0 {
		lw = -lw
	}
	return &StrokeProc{
		RHSProc:      TraceProc{rw / 2, d, JoinBevel},
		LHSProc:      TraceProc{lw / 2, d, JoinBevel},
		PostSideProc: nil,
		Width:        rw - lw,
		PointFunc:    PointCircle,
		CapStartFunc: cf,
		CapEndFunc:   cf,
	}
}
