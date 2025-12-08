package graphics2d

// StrandedProc takes a path and converts it into a number of parallel paths.
type StrandedProc struct {
	Traces []TraceProc
}

// NewStrandedProc returns a new stranded path processor.
func NewStrandedProc(n int, w float64) StrandedProc {
	if n < 2 {
		n = 2
	}
	if w < 0 {
		w = -w
	}
	dw := w / float64(n-1)
	w = -w / 2
	traces := make([]TraceProc, n)
	for i := 0; i < n; i++ {
		traces[i] = NewTraceProc(w)
		w += dw
	}
	return StrandedProc{traces}
}

// Process implements the PathProcessor interface.
func (sp StrandedProc) Process(p *Path) []*Path {
	paths := []*Path{}
	for _, tp := range sp.Traces {
		paths = append(paths, tp.Process(p)...)
	}
	return paths
}
