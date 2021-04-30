package graphics2d

// General purpose interface that takes a shape and turns it into
// a slice of shapes.

// ShapeProcessor defines the interface required for function passed to the Process function in Shape.
type ShapeProcessor interface {
	Process(s *Shape) []*Shape
}
