package graphics2d

import (
	"math"

	"github.com/jphsd/graphics2d/util"
)

// CapsProc contains a trio of shapes, one which will be placed using the start at the path, one at each step
// and the last, at the end. If either shape is nil, then it is skipped. The rotation flag indicates if the shapes
// should be rotated relative to the path's tangent at that point.
type CapsProc struct {
	Start  *Shape
	End    *Shape
	Mid    *Shape
	Rotate bool
}

// Process implements the PathProcessor interface.
func (cp CapsProc) Process(p *Path) []*Path {
	res := make([]*Path, 0, 2)
	parts := p.Parts()
	lp := len(parts)
	if cp.Rotate {
		if cp.Start != nil {
			t0 := util.DeCasteljau(parts[0], 0)
			xfm := CreateAffineTransform(t0[0], t0[1], 1, math.Atan2(t0[3], t0[2]))
			res = append(res, cp.Start.Transform(xfm).Paths()...)
		}
		if cp.Mid != nil {
			for i := 1; i < lp; i++ {
				t0 := util.DeCasteljau(parts[i], 0)
				xfm := CreateAffineTransform(t0[0], t0[1], 1, math.Atan2(t0[3], t0[2]))
				res = append(res, cp.Mid.Transform(xfm).Paths()...)
			}
		}
		// Only apply end cap to open paths
		if cp.End != nil && !p.Closed() {
			t1 := util.DeCasteljau(parts[lp-1], 1)
			xfm := CreateAffineTransform(t1[0], t1[1], 1, math.Atan2(t1[3], t1[2]))
			res = append(res, cp.End.Transform(xfm).Paths()...)
		}
		return res
	}

	if cp.Start != nil {
		part := parts[0]
		xfm := CreateAffineTransform(part[0][0], part[0][1], 1, 0)
		res = append(res, cp.Start.Transform(xfm).Paths()...)
	}
	if cp.Mid != nil {
		for i := 1; i < lp; i++ {
			part := parts[i]
			xfm := CreateAffineTransform(part[0][0], part[0][1], 1, 0)
			res = append(res, cp.Mid.Transform(xfm).Paths()...)
		}
	}
	if cp.End != nil && !p.Closed() {
		part := parts[lp-1]
		xfm := CreateAffineTransform(part[len(part)-1][0], part[len(part)-1][1], 1, 0)
		res = append(res, cp.End.Transform(xfm).Paths()...)
	}
	return res
}
