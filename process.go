package graphics2d

// General purpose interface that takes a path and turns it into
// a slice of paths. For example, a stroked outline of the path or breaks the
// path up into a series of dashes.

// PathProcessor defines the interface required for function passed to the Process function in Path.
type PathProcessor interface {
	Process(p *Path) []*Path
}
