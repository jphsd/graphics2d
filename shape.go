package graphics2d

import (
	"encoding/json"
	"fmt"
	"github.com/jphsd/graphics2d/util"
	"image"
	"math"
)

// Shape is a fillable collection of paths. For a path to be fillable,
// it must be closed, so paths added to the shape are forced closed on rendering.
type Shape struct {
	paths  []*Path
	bbox   [][]float64
	mask   *image.Alpha
	parent *Shape
}

// BoundingBox calculates a bounding box that the Shape is guaranteed to fit within.
func (s *Shape) BoundingBox() [][]float64 {
	if s.bbox == nil {
		var bb [][]float64
		for _, path := range s.paths {
			if bb == nil {
				bb = path.BoundingBox()
			} else {
				bbp := path.BoundingBox()
				bb = util.BoundingBox(bb[0], bb[1], bbp[0], bbp[1])
			}
		}
		s.bbox = bb
	}
	return s.bbox
}

// Bounds calculates the union of the bounds of the paths the shape contains.
func (s *Shape) Bounds() image.Rectangle {
	bb := s.BoundingBox()
	if bb == nil {
		return image.Rectangle{}
	}
	fx, fy := int(math.Floor(bb[0][0])), int(math.Floor(bb[0][1]))
	cx, cy := int(math.Ceil(bb[1][0])), int(math.Ceil(bb[1][1]))
	return image.Rectangle{image.Point{fx, fy}, image.Point{cx, cy}}
}

// Mask returns an Alpha image defined by the shape's bounds, containing the result
// of rendering the shape.
func (s *Shape) Mask() *image.Alpha {
	if s.mask != nil {
		return s.mask
	}
	srect := s.Bounds()
	s.mask = image.NewAlpha(srect)
	RenderShape(s.mask, s, image.Opaque)
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

// AddPaths adds paths to the shape.
func (s *Shape) AddPaths(paths ...*Path) {
	for _, p := range paths {
		if p == nil {
			continue
		}
		lp := p.Copy()
		if s.paths == nil {
			s.paths = make([]*Path, 1)
			s.paths[0] = lp
		} else {
			s.paths = append(s.paths, lp)
		}
	}
	s.bbox = nil
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

// Copy creates a new instance of this shape with a copy of its paths.
func (s *Shape) Copy() *Shape {
	paths := make([]*Path, len(s.paths))
	for i, path := range s.paths {
		paths[i] = path.Copy()
	}

	return &Shape{paths, nil, nil, s.parent}
}

// Transform applies a transform to all the paths in the shape
// and returns a new shape.
// JH deprecate this
func (s *Shape) Transform(xfm Transform) *Shape {
	return s.ProcessPaths(TransformProc{xfm})
}

// Transform applies a transform to all the paths in the shape

// Process applies a shape processor to the shape and
// returns a collection of new shapes.
func (s *Shape) Process(proc ShapeProcessor) []*Shape {
	shapes := proc.Process(s)
	// Fix parent
	for _, shape := range shapes {
		shape.parent = s
	}
	return shapes
}

// ProcessPaths applies a path processor to the shape and
// returns a new shape containing the processed paths.
func (s *Shape) ProcessPaths(proc PathProcessor) *Shape {
	np := make([]*Path, 0)
	for _, p := range s.paths {
		npaths := p.Process(proc)
		for _, pp := range npaths {
			np = append(np, pp)
		}
	}

	return &Shape{np, nil, nil, s}
}

// String converts a shape into a string.
func (s *Shape) String() string {
	str := fmt.Sprintf("SH %d ", len(s.paths))
	for _, path := range s.paths {
		str += path.String() + " "
	}
	return str
}

// PointInShape returns true if the point is contained within any path within the shape.
func (s *Shape) PointInShape(pt []float64) bool {
	for _, path := range s.paths {
		if path.PointInPath(pt) {
			return true
		}
	}
	return false
}

type shapeJSON struct {
	Paths []*Path
}

// MarshalJSON implements the encoding/json.Marshaler interface
func (s *Shape) MarshalJSON() ([]byte, error) {
	return json.Marshal(shapeJSON{s.paths})
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface
func (s *Shape) UnmarshalJSON(b []byte) error {
	var sj shapeJSON
	err := json.Unmarshal(b, &sj)
	if err != nil {
		return err
	}
	s.paths = sj.Paths

	// Reset everything else
	s.bbox = nil
	s.mask = nil
	s.parent = nil

	return nil
}
