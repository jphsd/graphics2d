package graphics2d

import (
	"github.com/jphsd/graphics2d/util"
	"math"
)

// RoundedProc replaces adjacent line segments in a path with line-arc-line where the radius of the
// arc is the minimum of Radius or the maximum allowable for the length of the shorter line segment.
// This ensures that the rounded corner doesn't end beyond the mid point of either line.
type RoundedProc struct {
	Radius float64
}

// Process implements the PathProcessor interface.
func (rp *RoundedProc) Process(p *Path) []*Path {
	parts := p.Parts()
	np := len(parts)
	if np < 2 {
		return []*Path{p}
	}

	res := [][][]float64{}
	for i, part := range parts {
		if len(part) != 2 {
			res = append(res, part)
			continue
		}
		if i < np-1 {
			if len(parts[i+1]) == 2 {
				nparts := rp.calcPieces(part[0], part[1], parts[i+1][1])
				res = append(res, [][]float64{part[0], nparts[0][0]})
				res = append(res, nparts...)
			} else {
				res = append(res, part)
			}
		} else {
			if !p.Closed() || len(parts[0]) != 2 {
				res = append(res, part)
				continue
			}
			// Path is closed and the first part is also a line
			nparts := rp.calcPieces(part[0], part[1], parts[0][1])
			res = append(res, [][]float64{part[0], nparts[0][0]})
			res = append(res, nparts...)
			lnp := len(nparts) - 1
			res[0][0] = nparts[lnp][len(nparts[lnp])-1]
		}
	}

	return []*Path{PartsToPath(res...)}
}

// Return p1-p2, p2-p3 intercepts and c, and final r and theta
func (rp *RoundedProc) calcPieces(p1, p2, p3 []float64) [][][]float64 {
	theta := util.AngleBetweenLines(p1, p2, p3, p2)
	neg := theta < 0
	if neg {
		theta = -theta
	}
	t2 := theta / 2
	tt2 := math.Tan(t2)

	// Check r is < min(p12, p23) / 2
	v1, v2 := util.Vec(p1, p2), util.Vec(p2, p3)
	d1, d2 := util.VecMag(v1), util.VecMag(v2)
	m1, m2 := d1/2, d2/2
	md := m1
	if m2 < m1 {
		md = m2
	}
	r := tt2 * md
	if r > rp.Radius {
		r = rp.Radius
	}

	// Find intersection of arc with p1-p2
	u1 := []float64{v1[0] / d1, v1[1] / d1}
	s := r / tt2
	i12 := []float64{p2[0] - s*u1[0], p2[1] - s*u1[1]}

	// Calc center
	c := []float64{i12[0], i12[1]}
	theta = math.Pi - theta
	n1 := []float64{u1[1], -u1[0]}
	if neg {
		n1[0], n1[1] = -n1[0], -n1[1]
	} else {
		theta = -theta
	}
	c = []float64{c[0] + r*n1[0], c[1] + r*n1[1]}

	// Calc offset
	a12 := math.Atan2(-n1[1], -n1[0])
	return MakeArcParts(c[0], c[1], r, a12, theta)
}
