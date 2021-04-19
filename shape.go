package graphics2d

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

// Shape is a fillable collection of paths. For a path to be fillable,
// it must be closed, so paths added to the shape are forced closed.
type Shape struct {
	paths  []*Path
	bounds image.Rectangle
	mask   *image.Alpha
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

// Mask returns an Alpha image, the size of the shape bounds, containing the result
// of rendering the shape.
func (s *Shape) Mask() *image.Alpha {
	if s.mask != nil {
		return s.mask
	}
	s.mask = image.NewAlpha(s.Bounds())
	min := s.Bounds().Min
	RenderShapeAlpha(s.mask, s, []float32{float32(-min.X), float32(-min.Y)}, draw.Over)
	return s.mask
}

// Contains returns true if the points are contained within the shape, false otherwise.
func (s *Shape) Contains(pts ...[]float64) bool {
	rect := s.Bounds()
	mask := s.Mask()
	ox, oy := rect.Min.X, rect.Min.Y
	mx, my := rect.Max.X, rect.Max.Y
	for _, pt := range pts {
		x := int(math.Floor(pt[0] + 0.5))
		y := int(math.Floor(pt[1] + 0.5))
		// Bounding box test
		if x < ox || x >= mx || y < oy || y >= my {
			return false
		}
		// Mask test
		if mask.AlphaAt(x, y).A < 128 {
			return false
		}
	}
	return true
}

// NewShape constructs a shape from the supplied paths.
func NewShape(paths ...*Path) *Shape {
	res := &Shape{}
	res.AddPaths(paths...)
	return res
}

// AddPaths adds paths to the shape and closes them if not already closed.
func (s *Shape) AddPaths(paths ...*Path) {
	for _, p := range paths {
		lp := p.Copy()
		lp.Close()
		if s.paths == nil {
			s.paths = make([]*Path, 1)
			s.paths[0] = lp
		} else {
			s.paths = append(s.paths, lp)
		}
	}
	s.bounds = image.Rectangle{}
	s.mask = nil
}

// AddShapes adds the paths from the supplied shapes to this shape.
func (s *Shape) AddShapes(shapes ...*Shape) {
	for _, shape := range shapes {
		s.AddPaths(shape.Paths()...)
	}
}

// Paths returns a shallow copy of the paths contained by this shape.
func (s *Shape) Paths() []*Path {
	return s.paths[:]
}

// Copy creates a new instance of this shape with a shallow copy of its paths.
func (s *Shape) Copy() *Shape {
	np := make([]*Path, len(s.paths))
	copy(np, s.paths)
	return &Shape{np, s.bounds, s.mask}
}

// Transform applies an affine transform to all the paths in the shape
// and returns a new shape.
func (s *Shape) Transform(xfm *Aff3) *Shape {
	np := make([]*Path, len(s.paths))
	for i, path := range s.paths {
		np[i] = path.Transform(xfm)
	}
	return &Shape{np, image.Rectangle{}, nil}
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

	return &Shape{np, image.Rectangle{}, nil}
}

// String converts a shape into a string.
func (s *Shape) String() string {
	str := fmt.Sprintf("SH %d ", len(s.paths))
	for _, path := range s.paths {
		str += path.String() + " "
	}
	return str
}
