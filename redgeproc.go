package graphics2d

import "math"

// RoundedEdgeProc is a path processor that replaces the parts of a path with an arc defined by
// the end points of the path and a third point normal to the part midpoint at either an absolute
// or relative (to the part length) distance from the midpoint. If Elip is set, then an elliptical arc
// of ry = d, rx = edge length / 2 is used instead.
type RoundedEdgeProc struct {
	Dist float64
	Abs  bool
	Elip bool
}

// Process implements the PathProcessor interface.
func (rp *RoundedEdgeProc) Process(p *Path) []*Path {
	parts := p.Parts()
	nparts := [][][]float64{}
	for _, part := range parts {
		cs, ce := part[0], part[len(part)-1]
		dx, dy := ce[0]-cs[0], ce[1]-cs[1]
		l := math.Hypot(dx, dy)
		mpx, mpy := cs[0]+dx/2, cs[1]+dy/2
		d := -rp.Dist
		if !rp.Abs {
			d *= l
		}
		nx, ny := -dy/l, dx/l
		var cpath *Path
		if rp.Elip {
			th := math.Atan2(dy, dx)
			ang := math.Pi
			var ry float64
			if d < 0 {
				ry = -d
			} else {
				ry = d
				ang = -ang
			}
			cp := []float64{mpx, mpy}
			cpath = EllipticalArc(cp, l/2, ry, math.Pi, ang, 0, ArcOpen)
			xfm := RotateAbout(th, mpx, mpy)
			cpath = cpath.Process(&TransformProc{xfm})[0]
		} else {
			cm := []float64{mpx + nx*d, mpy + ny*d}
			cpath = ArcFromPoints(cs, cm, ce, ArcOpen)
		}
		nparts = append(nparts, cpath.Parts()...)
	}
	path := PartsToPath(nparts...)
	if p.Closed() {
		path.Close()
	}
	return []*Path{path}
}
