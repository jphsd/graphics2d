package graphics2d

import "math"

// RoundedEdgeProc is a path processor that replaces the parts of a path with an arc defined by
// the end points of the path and a third point normal to the part midpoint at either an absolute
// or relative (to the part length) distance from the midpoint.
type RoundedEdgeProc struct {
	Dist float64
	Abs  bool
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
		cm := []float64{mpx + nx*d, mpy + ny*d}
		cparts := ArcFromPoints(cs, cm, ce, ArcOpen).Parts()
		nparts = append(nparts, cparts...)
	}
	path := PartsToPath(nparts...)
	if p.Closed() {
		path.Close()
	}
	return []*Path{path}
}
