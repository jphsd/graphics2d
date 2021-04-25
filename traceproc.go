package graphics2d

// Constant width path tracer. Traces a path at a normal distance of width from the path.
// Join types - round, bevel [default], miter

// TraceProc defines the width and join types of the trace. The gap between two adjacent
// steps must be greater than MinGap for the join function to be called.
type TraceProc struct {
	Width    float64
	MinGap   float64
	Flatten  float64
	JoinFunc func([][]float64, []float64, [][]float64) [][][]float64
}

// NewTraceProc creates a stroke path processor with width w, the bevel join and butt cap types.
func NewTraceProc(w float64) *TraceProc {
	return &TraceProc{w, 1, 0.5, JoinBevel} // 10 degrees
}

// Process implements the PathProcessor interface.
func (tp *TraceProc) Process(p *Path) []*Path {
	path, _ := PartsToPath(tp.ProcessParts(p)...)
	if path == nil {
		return []*Path{}
	}
	if p.Closed() {
		path.Close()
	}

	if tp.Width < 0 {
		path = path.Reverse()
	}

	return []*Path{path}
}

// ProcessParts returns the processed path as a slice of parts, rather a path so other path
// processors don't have to round trip path -> parts -> path -> parts (e.g. StrokeProc).
func (tp *TraceProc) ProcessParts(p *Path) [][][]float64 {
	// A point isn't traceable.
	if len(p.Steps()) == 1 {
		return [][][]float64{}
	}

	w := tp.Width
	if w < 0 {
		w = -w
		p = p.Reverse()
	}

	// Preprocess curves into safe forms.
	p = p.Simplify()
	parts := p.Parts()
	n := len(parts)
	tangs := p.Tangents()

	// Convert tangents to scaled RHS normals
	norms := make([][][]float64, n)
	for i := 0; i < n; i++ {
		norms[i] = make([][]float64, 2)
		norms[i][0] = []float64{w * tangs[i][0][1], -w * tangs[i][0][0]}
		norms[i][1] = []float64{w * tangs[i][1][1], -w * tangs[i][1][0]}
	}

	// Calculate the path by LineTransforming the parts and handling the joins
	rhs := make([][][]float64, n)
	for i := 0; i < n; i++ {
		part := parts[i]
		ln := len(part) - 1
		offs := norms[i]
		xfm := LineTransform(part[0][0], part[0][1],
			part[ln][0], part[ln][1],
			part[0][0]+offs[0][0], part[0][1]+offs[0][1],
			part[ln][0]+offs[1][0], part[ln][1]+offs[1][1])

		rhs[i] = xfm.Apply(part...)
	}

	// Compute the joins
	nrhs := make([][][]float64, 0, 2*n)
	nrhs = append(nrhs, rhs[0])
	for i := 1; i < n; i++ {
		last := nrhs[len(nrhs)-1]
		// Check for knot first
		npt := PartsIntersection(last, rhs[i], tp.Flatten)
		if npt != nil {
			// Tweak the end of nrhs[$] and start of rhs[i]
			// Not strictly correct - should really figure out the t value for
			// the point and then split part at t value to preserve the part's cp.
			last[len(last)-1] = npt
			rhs[i][0] = npt
		} else if !tp.checkGap(last[len(last)-1], rhs[i][0]) {
			nrhs = append(nrhs, tp.JoinFunc(last, parts[i][0], rhs[i])...)
		} else {
			// No join, so nudge end to start of rhs[i]
			last[len(last)-1] = rhs[i][0]
		}
		nrhs = append(nrhs, rhs[i])
	}

	if p.Closed() {
		// Join the end points
		last := nrhs[len(nrhs)-1]
		npt := PartsIntersection(last, nrhs[0], tp.Flatten)
		if npt != nil {
			last[len(last)-1] = npt
			nrhs[0][0] = npt
		} else if !tp.checkGap(last[len(last)-1], nrhs[0][0]) {
			nrhs = append(nrhs, tp.JoinFunc(last, parts[0][0], nrhs[0])...)
		} else {
			last[len(last)-1] = nrhs[0][0]
		}
	}

	return nrhs
}

func (tp *TraceProc) checkGap(p1, p2 []float64) bool {
	d2 := tp.MinGap * tp.MinGap
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]
	return d2 > dx*dx+dy*dy
}
