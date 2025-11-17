package graphics2d

import (
	"math"
	"math/rand"

	"github.com/jphsd/graphics2d/util"
)

// PointRot specifies the type of shape rotation
type PointRot int

const (
	// RotFixed no rotation
	RotFixed PointRot = iota
	// RotRelative rotation relative to the tangent of the path step
	RotRelative
	// RotRandom rotation is randomized
	RotRandom
)

// PointsProc contains a slice of shapes, one of which will be placed at the start of each step in the
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
		if partLen(part) < 0.0001 {
			// Skip 0 length parts which cause tangent issues
			continue
		}
		if pp.Points[cp] != nil {
			var xfm *Aff3
			switch pp.Rotate {
			default:
				fallthrough
			case RotFixed:
				xfm = CreateAffineTransform(part[0][0], part[0][1], 1, 0)
			case RotRelative:
				t0 := util.DeCasteljau(part, 0)
				ang := math.Atan2(t0[3], t0[2])
				xfm = CreateAffineTransform(t0[0], t0[1], 1, ang)
			case RotRandom:
				t0 := util.DeCasteljau(part, 0)
				xfm = CreateAffineTransform(t0[0], t0[1], 1, rand.Float64()*TwoPi)
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
			xfm = CreateAffineTransform(part[0][0], part[0][1], 1, 0)
		case RotRelative:
			t0 := util.DeCasteljau(part, 1)
			xfm = CreateAffineTransform(t0[0], t0[1], 1, math.Atan2(t0[3], t0[2]))
		case RotRandom:
			t0 := util.DeCasteljau(part, 1)
			xfm = CreateAffineTransform(t0[0], t0[1], 1, rand.Float64()*TwoPi)
		}
		res = append(res, pp.Points[cp].Transform(xfm).Paths()...)
	}
	return res
}

// ShapesProc contains a slice of shapes, which will be placed sequentially along the path,
// starting at the beginning and spaced there after by the spacing value, and at the path end,
// if not closed. If any shape is nil, then it is skipped. The rotation flag indicates if the
// shapes should be rotated relative to the path's tangent at that point.
type ShapesProc struct {
	Comp   *CompoundProc
	Shapes *PointsProc
}

// NewShapesProc creates a new shapes path processor with the supplied shapes, spacing and rotation flag.
func NewShapesProc(shapes []*Shape, spacing float64, rot PointRot) ShapesProc {
	pattern := []float64{spacing, spacing}
	spaces := NewSnipProc(2, pattern, 0)
	n := len(shapes)
	d := int(math.Floor(spacing + 0.5))
	nn := d * n
	nshapes := make([]*Shape, nn)
	for i := range n {
		nshapes[i*d] = shapes[i]
		// remaining d-1 slots are left nil
	}
	// Assumption - this is taking place in image space so pixel level sampling should
	// be sufficient.
	comp := NewCompoundProc(NewMunchProc(1), spaces)
	comp.Concatenate = true
	return ShapesProc{comp, NewPointsProc(nshapes, rot)}
}

// Process implements the PathProcessor interface.
func (sp ShapesProc) Process(p *Path) []*Path {
	path := p.Process(sp.Comp)[0]

	return path.Process(sp.Shapes)
}

func partLen(part [][]float64) float64 {
	p1 := part[0]
	p2 := part[len(part)-1]
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	return math.Hypot(dx, dy)
}
