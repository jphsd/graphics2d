package graphics2d

import "math/rand"

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
	shapes := make([]*Shape, np)
	for i, path := range paths {
		shapes[i] = NewShape(path)
	}

	return shapes
}

// BucketProc agregates paths into N shapes using the specificed style.
type BucketProc struct {
	N     int
	Style BucketStyle
}

type BucketStyle int

const (
	Chunk BucketStyle = iota
	RoundRobin
	Random
)

// Process implements the ShapeProcessor interface.
func (bp *BucketProc) Process(s *Shape) []*Shape {
	shapes := make([]*Shape, bp.N)
	paths := s.Paths()
	np := len(paths)
	b := 0
	switch bp.Style {
	case Chunk:
		npb := np / bp.N
		if npb < 1 {
			npb = 1
		}
		for i, path := range paths {
			shapes[b].AddPaths(path)
			if i > 0 && i%npb == 0 {
				b++
			}
		}
	case RoundRobin:
		for _, path := range paths {
			shapes[b].AddPaths(path)
			b++
			if b == bp.N {
				b = 0
			}
		}
	case Random:
		for _, path := range paths {
			shapes[rand.Intn(bp.N)].AddPaths(path)
		}
	}

	return shapes
}
