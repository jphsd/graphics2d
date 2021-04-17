package graphics2d

// General purpose interface that takes a path and turns it into
// a slice of paths. For example, a stroked outline of the path or breaks the
// path up into a series of dashes.

// PathProcessor defines the interface required for function passed to the Process function in Path.
type PathProcessor interface {
	Process(p *Path) []*Path
}

// CompoundProcessor applies a collection of PathProcessors to a path.
type CompoundProc struct {
	Procs []PathProcessor
}

// NewCompoundProc creates a new CompundProcessor with the supplied path processors.
func NewCompoundProc(procs ...PathProcessor) *CompoundProc {
	return &CompoundProc{procs}
}

// Process implements the PathProcessor interface.
func (cp *CompoundProc) Process(p *Path) []*Path {
	paths := []*Path{p}
	if len(cp.Procs) == 0 {
		return paths
	}

	for _, proc := range cp.Procs {
		npaths := []*Path{}
		for _, path := range paths {
			npaths = append(npaths, proc.Process(path)...)
		}
		paths = npaths
	}

	return paths
}

// PathProcessor wrappers for Flatten, Simplify and Linepath functions

// Flatten contains the minimum required distance to the controa points.
type FlattenProc struct {
	D float64
}

// Process implements the PathProcessor interface.
func (fp *FlattenProc) Process(p *Path) []*Path {
	path := p.Flatten(fp.D)
	return []*Path{path}
}

type LineProc struct{}

// Process implements the PathProcessor interface.
func (lp *LineProc) Process(p *Path) []*Path {
	path := p.Line()
	return []*Path{path}
}

type SimplifyProc struct{}

// Process implements the PathProcessor interface.
func (sp *SimplifyProc) Process(p *Path) []*Path {
	path := p.Simplify()
	return []*Path{path}
}
