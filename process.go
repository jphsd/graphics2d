package graphics2d

// General purpose interface that takes a path and turns it into
// a slice of paths. For example, a stroked outline of the path or breaks the
// path up into a series of dashes.

// PathProcessor defines the interface required for function passed to the Process function in Path.
type PathProcessor interface {
	Process(p *Path) []*Path
}

// CompoundProcess applies a collection of PathProcessors to a path.
type CompoundProcessor struct {
	Procs []PathProcessor
}

// NewCompoundProcessor creates a new CompundProcessor with the supplied path processors.
func NewCompoundProcessor(procs []PathProcessor) *CompoundProcessor {
	return &CompoundProcessor{procs}
}

// Process implements the PathProcessor interface.
func (cp *CompoundProcessor) Process(p *Path) []*Path {
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
