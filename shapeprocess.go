package graphics2d

// General purpose interface that takes a shape and turns it into
// a slice of shapes.

// ShapeProcessor defines the interface required for function passed to the Process function in Shape.
type ShapeProcessor interface {
	Process(s *Shape) []*Shape
}

// PathsProc converts each path in a shape into its own shape.
type PathsProc struct{}

// Process implements the ShapeProcessor interface.
func (pp *PathsProc) Process(s *Shape) []*Shape {
	paths := s.Paths()
	np := len(paths)
	shapes := make([]*shapes, np)
	for i, path := range paths {
		shapes[i] = NewShape(path)
	}

	return shapes
}
