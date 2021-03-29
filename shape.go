package graphics2d

import (
	"fmt"
	"image"
)

// Shape is a fillable collection of paths. For a path to be fillable,
// it must be closed, so paths added to the shape are forced closed.
type Shape struct {
	paths  []*Path
	bounds image.Rectangle
}

// Bounds calculates the union of the bounds of the paths the shape contains.
func (s *Shape) Bounds() image.Rectangle {
	if s.bounds.Empty() && s.paths != nil {
		rect := s.paths[0].Bounds()
		for i := 1; i < len(s.paths); i++ {
			rect = rect.Union(s.paths[i].Bounds())
		}
		s.bounds = rect
	}
	return s.bounds
}

// AddPath adds a path to the shape and closes it if not already closed.
func (s *Shape) AddPath(p *Path) {
	lp := p.Copy()
	lp.Close()
	if s.paths == nil {
		s.paths = make([]*Path, 1)
		s.paths[0] = lp
	} else {
		s.paths = append(s.paths, lp)
	}
	s.bounds = image.Rectangle{}
}

// AddPaths adds a collection of paths to the shape.
func (s *Shape) AddPaths(paths []*Path) {
	for _, p := range paths {
		s.AddPath(p)
	}
}

// AddShape adds the paths from the supplied shape to this shape.
func (s *Shape) AddShape(shape *Shape) {
	s.AddPaths(shape.Paths())
}

// Paths returns a shallow copy of the paths contained by this shape.
func (s *Shape) Paths() []*Path {
	return s.paths[:]
}

// Copy creates a new instance of this shape with a shallow copy of its paths.
func (s *Shape) Copy() *Shape {
	np := make([]*Path, len(s.paths))
	copy(np, s.paths)
	return &Shape{np, s.bounds}
}

// Transform applies an affine transform to all the paths in the shape
// and returns a new shape.
func (s *Shape) Transform(xfm *Aff3) *Shape {
	np := make([]*Path, len(s.paths))
	for i, path := range s.paths {
		np[i] = path.Transform(xfm)
	}
	return &Shape{np, image.Rectangle{}}
}

// Process applies a processor to the shape and
// returns a new shape containing the processed paths.
func (s *Shape) Process(proc PathProcessor) *Shape {
	np := make([]*Path, 0)
	for _, p := range s.paths {
		npaths := p.Process(proc)
		for _, pp := range npaths {
			pp.Close()
			np = append(np, pp)
		}
	}

	return &Shape{np, image.Rectangle{}}
}

// String converts a shape into a string.
func (s *Shape) String() string {
	str := fmt.Sprintf("SH %d ", len(s.paths))
	for _, path := range s.paths {
		str += path.String() + " "
	}
	return str
}
