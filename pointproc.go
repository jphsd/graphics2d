package graphics2d

import (
	"math"
	"math/rand"

	. "github.com/jphsd/graphics2d/util"
)

// PointRot specifies the type of shape rotation
type PointRot int

const (
	// RotFixed, no rotation
	RotFixed PointRot = iota
	// RotRelative, rotation relative to the tangent of the path step
	RotRelative
	// RotRandom, rotation is randomized
	RotRandom
)

// PointsProc contains a slice of shapes, one of which will be placed using at the start of each step in the
// path and at the path end, if not closed. If any shape is nil, then it is skipped. The rotation flag
// indicates if the shapes should be rotated relative to the path's tangent at that point.
type PointsProc struct {
	Points []*Shape
	Rotate PointRot
}

// NewPointsProc creates a new points path processor with the supplied shapes and rotation flag.
func NewPointsProc(shapes []*Shape, rot PointRot) *PointsProc {
	return &PointsProc{shapes, rot}
}

// Process implements the PathProcessor interface.
func (pp *PointsProc) Process(p *Path) []*Path {
	parts := p.Parts()
	n := len(parts)
	if !p.Closed() {
		n++
	}
	res := make([]*Path, 0, n)
	ns := len(pp.Points)
	cp := 0
	for _, part := range parts {
		if pp.Points[cp] != nil {
			var xfm *Aff3
			switch pp.Rotate {
			case RotFixed:
				xfm = CreateTransform(part[0][0], part[0][1], 1, 0)
			case RotRelative:
				t0 := DeCasteljau(part, 0)
				xfm = CreateTransform(t0[0], t0[1], 1, math.Atan2(t0[3], t0[2]))
			case RotRandom:
				t0 := DeCasteljau(part, 0)
				xfm = CreateTransform(t0[0], t0[1], 1, rand.Float64()*math.Pi*2)
			}
			res = append(res, pp.Points[cp].Transform(xfm).Paths()...)
		}
		if cp++; cp == ns {
			cp = 0
		}
	}
	// Only apply end shape to open paths
	if pp.Points[cp] != nil && !p.Closed() {
		var xfm *Aff3
		part := parts[len(parts)-1]
		switch pp.Rotate {
		case RotFixed:
			xfm = CreateTransform(part[0][0], part[0][1], 1, 0)
		case RotRelative:
			t0 := DeCasteljau(part, 1)
			xfm = CreateTransform(t0[0], t0[1], 1, math.Atan2(t0[3], t0[2]))
		case RotRandom:
			t0 := DeCasteljau(part, 1)
			xfm = CreateTransform(t0[0], t0[1], 1, rand.Float64()*math.Pi*2)
		}
		res = append(res, pp.Points[cp].Transform(xfm).Paths()...)
	}
	return res
}

// ShapesProc contains a slice of shapes, which will be placed sequentially using along the path,
// starting at the beginning and spaced there after by the spacing value, and at the path end,
// if not closed. If any shape is nil, then it is skipped. The rotation flag indicates if the
// shapes should be rotated relative to the path's tangent at that point.
type ShapesProc struct {
	Munch  *MunchProc
	Spaces *SnipProc
	Shapes *PointsProc
}

// NewShapesProc creates a new shapes path processor with the supplied shapes, spacing and rotation flag.
func NewShapesProc(shapes []*Shape, spacing float64, rot PointRot) *ShapesProc {
	pattern := []float64{spacing, spacing}
	spaces := NewSnipProc(2, pattern, 0)
	n := len(shapes)
	d := int(math.Floor(spacing + 0.5))
	nn := d * n
	nshapes := make([]*Shape, nn)
	for i := 0; i < n; i++ {
		nshapes[i*d] = shapes[i]
	}
	return &ShapesProc{NewMunchProc(1), spaces, NewPointsProc(nshapes, rot)}
}

// Process implements the PathProcessor interface.
func (sp *ShapesProc) Process(p *Path) []*Path {
	// Break up the path into pieces smaller than the spacingg so the tangents are
	// correct - Spaces has been grown with nil shapes to accommodate the extra steps.
	paths := p.Process(sp.Munch)
	path, _ := ConcatenatePaths(paths...)
	paths = path.Process(sp.Spaces)
	path, _ = ConcatenatePaths(paths...)
	return path.Process(sp.Shapes)
}
