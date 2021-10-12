package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// CapsProc contains a pair of shapes, one which will be placed using the start at the path and the other
// at the end. If either shape is nil, then it is skipped. The rotation flag indicates if the shapes
// should be rotated relative to the path's tangent at that point.
type CapsProc struct {
	Caps   []*Shape
	Rotate bool
}

// NewCapsProc creates a new cap path processor with the supplied shapes and rotation flag.
func NewCapsProc(start, end *Shape, rot bool) *CapsProc {
	return &CapsProc{[]*Shape{start, end}, rot}
}

// Process implements the PathProcessor interface.
func (cp *CapsProc) Process(p *Path) []*Path {
	res := make([]*Path, 0, 2)
	parts := p.Parts()
	if cp.Rotate {
		if cp.Caps[0] != nil {
			t0 := util.DeCasteljau(parts[0], 0)
			xfm := CreateTransform(t0[0], t0[1], 1, math.Atan2(t0[3], t0[2]))
			res = append(res, cp.Caps[0].Transform(xfm).Paths()...)
		}
		// Only apply end cap to open paths
		if cp.Caps[1] != nil && !p.Closed() {
			t1 := util.DeCasteljau(parts[len(parts)-1], 1)
			xfm := CreateTransform(t1[0], t1[1], 1, math.Atan2(t1[3], t1[2]))
			res = append(res, cp.Caps[1].Transform(xfm).Paths()...)
		}
		return res
	}

	if cp.Caps[0] != nil {
		part := parts[0]
		xfm := CreateTransform(part[0][0], part[0][1], 1, 0)
		res = append(res, cp.Caps[0].Transform(xfm).Paths()...)
	}
	if cp.Caps[1] != nil && !p.Closed() {
		part := parts[len(parts)-1]
		xfm := CreateTransform(part[len(part)-1][0], part[len(part)-1][1], 1, 0)
		res = append(res, cp.Caps[1].Transform(xfm).Paths()...)
	}
	return res
}
