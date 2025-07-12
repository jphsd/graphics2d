package graphics2d

// TransformProc produces a new path by applying the supplied transform to the path.
type TransformProc struct {
	Transform Transform
}

// Process implements the PathProcessor interface.
func (tp *TransformProc) Process(p *Path) []*Path {
	steps := p.Steps()
	for i, step := range steps {
		steps[i] = tp.Transform.Apply(step...)
	}

	path := NewPath(steps[0][0])
	for _, step := range steps {
		path.AddStep(step...)
	}

	if p.Closed() {
		path.Close()
	}

	return []*Path{path}
}
