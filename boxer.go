package graphics2d

import "math"

// BoxerProc is a path processor that converts a path into a set of boxes along the path with the
// specified width. If the offset is 0 then the box is centered on the path.
type BoxerProc struct {
	Width float64 // Width of box
	Offs  float64 // Offset from path of box
	Flat  float64 // Flatten value
}

// NewBoxerProc returns a new BoxerProc path processor.
func NewBoxerProc(width, offs float64) *BoxerProc {
	return &BoxerProc{width, offs, RenderFlatten}
}

// Process implements the PathProcessor interface.
func (bp *BoxerProc) Process(p *Path) []*Path {
	paths := []*Path{}
	hw := bp.Width / 2
	for _, part := range p.Flatten(bp.Flat).Parts() {
		paths = append(paths, Polygon(box(part[0], part[1], hw, bp.Offs)...))
	}

	return paths
}

// Convert start and end points into a set of box corners.
func box(p1, p2 []float64, hw, offs float64) [][]float64 {
	dx, dy := p2[0]-p1[0], p2[1]-p1[1]
	nx, ny := -dy, dx
	d := math.Sqrt(nx*nx + ny*ny)
	nx /= d
	ny /= d
	l, r := hw-offs, hw+offs
	lnx, lny := l*nx, l*ny
	rnx, rny := r*nx, r*ny
	return [][]float64{
		{p1[0] + rnx, p1[1] + rny},
		{p2[0] + rnx, p2[1] + rny},
		{p2[0] - lnx, p2[1] - lny},
		{p1[0] - lnx, p1[1] - lny},
	}
}
