package graphics2d

// General purpose interface that takes a path and turns it into
// a slice of paths. For example, a stroked outline of the path or breaks the
// path up into a series of dashes.

// PathProcessor defines the interface required for function passed to the Process function in Path.
type PathProcessor interface {
	Process(p *Path) []*Path
}

// CompoundProc applies a collection of PathProcessors to a path.
type CompoundProc struct {
	Procs       []PathProcessor
	Concatenate bool // If set, concatenate processed paths into a single path after every processor
}

// NewCompoundProc creates a new CompoundProcessor with the supplied path processors.
func NewCompoundProc(procs ...PathProcessor) CompoundProc {
	return CompoundProc{procs, false}
}

// Process implements the PathProcessor interface.
func (cp CompoundProc) Process(p *Path) []*Path {
	paths := []*Path{p}
	if len(cp.Procs) == 0 {
		return paths
	}

	for _, proc := range cp.Procs {
		npaths := []*Path{}
		for _, path := range paths {
			npaths = append(npaths, proc.Process(path)...)
		}
		if cp.Concatenate {
			path, _ := ConcatenatePaths(npaths...)
			paths = []*Path{path}
		} else {
			paths = npaths
		}
	}

	if cp.Concatenate && p.Closed() {
		paths[0].Close()
	}

	return paths
}

// PathProcessor wrappers for Flatten, Simplify, and Reverse path functions

// FlattenProc is a wrapper around Path.Flatten() and contains the minimum required
// distance to the control points.
type FlattenProc struct {
	Flatten float64
}

// Process implements the PathProcessor interface.
func (fp FlattenProc) Process(p *Path) []*Path {
	return []*Path{p.Flatten(fp.Flatten)}
}

// ReverseProc replaces a path with its reverse.
type ReverseProc struct{}

// Process implements the PathProcessor interface.
func (rp ReverseProc) Process(p *Path) []*Path {
	return []*Path{p.Reverse()}
}

// SimplifyProc is a wrpper around Path.Simplify().
type SimplifyProc struct{}

// Process implements the PathProcessor interface.
func (sp SimplifyProc) Process(p *Path) []*Path {
	return []*Path{p.Simplify()}
}

// StepsProc converts each path step into its own path.
type StepsProc struct{}

// Process implements the PathProcessor interface.
func (sp StepsProc) Process(p *Path) []*Path {
	parts := p.Parts()
	np := len(parts)
	paths := make([]*Path, np)
	for i, part := range parts {
		paths[i] = PartsToPath(part)
	}

	return paths
}
