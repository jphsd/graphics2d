package graphics2d

// ShapeProc creates a new shape for every path in the shape.
type ShapeProc struct{}

// Process implements the ShapeProcessor interface.
func (sp *ShapeProc) Process(s *Shape) []*Shape {
	paths := s.Paths()
	shapes := make([]*Shape, len(paths))
	for i := range len(paths) {
		shapes[i] = NewShape(paths[i])
	}

	return shapes
}
